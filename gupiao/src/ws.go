package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	OneHand   = 100
	WarnHL    = 1
	WarnRatio = 200
	LB        = 300
)
const wsc = 0

var checkSecs = []int64{3, 10, 30, 60, 120, 300}
var checkCnts = []int{10, 20, 30, 50, 100}

type empty struct {
}

var reqid int32 = 1

type HL struct {
	Min float64
	Max float64
}

// 最高价和最低价
var mId2HL map[string]*HL

// 开盘时间
var mId2BaseData map[string]StatisticType

// 基本信息
var mId2ConstInfo map[string]StaticType

// 最近买单
var mId2Tick map[string][]TickType

// 最后价格
var mId2Dyna map[string][]DynaType

// 量比
var mId2LB map[string]*VRa

// 监听者
var listerMx sync.Mutex
var mId2Listener map[string]map[string]empty
var mId2Post map[string]empty
var mId2Time map[string]int64
var mId2cnt map[string]int64

var wscon []*websocket.Conn

func postJson(id int, data interface{}) {
	wscon[id].WriteJSON(data)
}
func hx(id string) int {
	x, _ := strconv.Atoi(id)
	return x % (wsc + 1)
}

func PostSTATISTICS(id, dq string) {

	dataJson := Data_json{
		SubType:     "SUBON",
		Inst:        id,
		Market:      dq,
		ServiceType: "STATISTICS",
		ReqID:       int(atomic.AddInt32(&reqid, 1)),
	}
	postJson(hx(id), dataJson)

}
func PostTick(id, dq string) {
	dataJson := Data_json{
		SubType:     "SUBON",
		Inst:        id,
		Market:      dq,
		ServiceType: "TICK",
		ReqID:       int(atomic.AddInt32(&reqid, 1)),
	}
	postJson(hx(id), dataJson)

}
func PostStatic(id, dq string) {

	dataJson := Data_json{
		SubType:     "SUBON",
		Inst:        id,
		Market:      dq,
		ServiceType: "STATIC",
		ReqID:       int(atomic.AddInt32(&reqid, 1)),
	}
	postJson(hx(id), dataJson)

}
func PostDyna(id, dq string) {

	dataJson := Data_json{
		SubType:     "SUBON",
		Inst:        id,
		Market:      dq,
		ServiceType: "DYNA",
		ReqID:       int(atomic.AddInt32(&reqid, 1)),
	}
	postJson(hx(id), dataJson)

}

func Post(id, dq string) {
	//当前没有post过
	if _, ok := mId2Post[id]; !ok {
		PostStatic(id, dq)
		PostSTATISTICS(id, dq)
		PostTick(id, dq)
		PostDyna(id, dq)
		mId2Post[id] = empty{}
	}
	listerMx.Lock()
	if _, ok := mId2Listener[id]; !ok {
		mId2Listener[id] = map[string]empty{}
	}
	listerMx.Unlock()

}
func PostAndListen(name, id, dq string) {
	//当前没有post过
	Post(id, dq)
	listerMx.Lock()

	mId2Listener[id][name] = empty{}
	listerMx.Unlock()

}

func getRa(cur, base float64) float64 {
	return (cur - base) / base * 100
}

func Ping(conn *websocket.Conn) {
	ping := PingType{ServiceType: "ping"}
	ticker := time.NewTicker(60 * time.Second)
	select {
	case <-ticker.C:
		{
			conn.WriteJSON(ping)
		}

	}
}

// 股票异动
func handleTick(r dataRes) {
	muban := ""
	v := []int{}
	str := []string{}

	ramsg := ""
	bra := false
	if info, ok := mId2ConstInfo[r.Inst]; ok {
		muban = info.InstrumentName + "\n"
		ramsg = fmt.Sprintf("%s  %s  交易量变大\n", r.Inst, info.InstrumentName)
	}
	b := false
	notify := false
	ch := "↑"
	lastPirce := float64(0)
	ra, _ok := mId2LB[r.Inst]
	if !_ok {
		ra = &VRa{sum: 0, n: 0}
		mId2LB[r.Inst] = ra
	}
	for _, x := range r.QuoteData.TickData {
		notify = true
		mId2Tick[r.Inst] = append(mId2Tick[r.Inst], x)
		tprice := x.Price * float64(x.Volume)

		//fmt.Println(tprice / ra.GetAvg())
		if mId2Time[r.Inst]+6 < x.Time && tprice/ra.GetAvg() >= WarnRatio {
			fmt.Println("log1 ", tprice/ra.GetAvg(), tprice, ra.GetAvg(), r.Inst)

			bra = true
			mId2Time[r.Inst] = x.Time
		}
		ra.Push(VRaInnner{val: tprice, t: x.Time})

		if x.Price >= lastPirce {
			ch = "↑"
		} else {
			ch = "↓"
		}
		lastPirce = x.Price
		if x.Volume > 200*OneHand {
			base := x.Price
			if sts, ok := mId2BaseData[r.Inst]; ok {
				base = sts.PreClosePrice
			}
			b = true
			//fmt.Println("???", x.Price, base, getRa(x.Price, base))
			if getRa(x.Price, base) < 22 {
				str = append(str, fmt.Sprintf("%.02f%%   %.02f%s   %d\n", getRa(x.Price, base), x.Price, ch, x.Volume/OneHand))
				v = append(v, x.Volume)
			} else {
				fmt.Println("err,", x.Price, base, getRa(x.Price, base))
			}
		}
	}

	if b {
		n := len(v)

		listerMx.Lock()
		for name, _ := range mId2Listener[r.Inst] {
			needlen := getFollow(name).FollowsId[r.Inst].WarnMsg
			smsg := muban
			b = false
			for i := 0; i < n; i++ {
				if v[i] >= needlen {
					smsg += str[i]
					b = true
				}
			}
			if b {
				SendMsg(name, smsg)
			}
		}

		listerMx.Unlock()
	}
	if notify {
		run := true
		for _, sec := range checkSecs {
			run = checkUnActionByTime(r.Inst, sec, run)
		}
		run = true
		for _, cnt := range checkCnts {
			run = checkUnActionByCount(r.Inst, cnt, run)
		}

	}
	if bra {

		listerMx.Lock()
		for name, _ := range mId2Listener[r.Inst] {
			SendMsg(name, ramsg)
		}

		listerMx.Unlock()
		ddMsg <- ramsg
	}
}
func SendMsg2Listen(inst, msg string) {

	listerMx.Lock()
	for name, _ := range mId2Listener[inst] {
		SendMsg(name, msg)
	}

	listerMx.Unlock()
}
func checkUnActionByTime(id string, sec int64, run bool) bool {
	if run == false {
		return false
	}
	//算法 记录最高点和最低点 如果当前点不是最新的最高点和最低点就不上报
	low := float64(60000)
	lidx := -1
	hight := float64(0) //A股不会再有6000点
	hidx := -1
	n := len(mId2Tick[id])
	if n <= 0 {
		return false
	}
	endsec := mId2Tick[id][n-1].Time
	for i := n - 1; i >= 0; i-- {
		tmp := mId2Tick[id][i]
		if endsec-tmp.Time <= sec {
			if tmp.Price >= hight {
				hight = tmp.Price
				hidx = i
			}
			if tmp.Price <= low {
				low = tmp.Price
				lidx = i
			}
		} else {
			break
		}
	}
	if (lidx == n-1 || hidx == n-1) && lidx != hidx {

		muban := ""
		if info, ok := mId2ConstInfo[id]; ok {
			muban = info.InstrumentName
		}
		if sts, ok := mId2BaseData[id]; ok {
			base := sts.PreClosePrice
			ra1 := getRa(hight, base)
			ra2 := getRa(low, base)
			diff := ra1 - ra2
			//fmt.Println("...", hidx, low, diff, base)
			if diff >= WarnRatio && diff < 20 {
				ch := "↑↑↑"
				if lidx == n-1 {
					ch = "↓↓↓"
				}
				SendMsg2Listen(id, fmt.Sprintf("%s 在%d秒异动%s\n%.2f%%  %.2f%%  %.2f%% \n", muban, sec, ch, getRa(low, base), getRa(hight, base), diff))
				fmt.Printf("%s 在%d秒异动 %s %.2f%%  %.2f%%  %.2f%%  %s %.2f %.2f %v\n", muban, sec, ch, getRa(low, base), getRa(hight, base), diff, id, low, hight, base)
				return false
			}
		}
	}
	return true
}
func checkUnActionByCount(id string, cnt int, run bool) bool {
	if run == false || mId2cnt[id]+60 > time.Now().Unix() {
		return false
	}
	tcnt := cnt
	//算法 记录最高点和最低点 如果当前点不是最新的最高点和最低点就不上报
	low := float64(60000)
	lidx := -1
	hight := float64(0) //A股不会再有6000点
	hidx := -1
	n := len(mId2Tick[id])
	if n <= 0 {
		return false
	}
	for i := n - 1; i >= 0 && cnt > 0; i-- {
		cnt--
		tmp := mId2Tick[id][i]
		if tmp.Price >= hight {
			hight = tmp.Price
			hidx = i
		}
		if tmp.Price <= low {
			low = tmp.Price
			lidx = i
		}
	}
	if (lidx == n-1 || hidx == n-1) && lidx != hidx {

		muban := ""
		if info, ok := mId2ConstInfo[id]; ok {
			muban = info.InstrumentName
		}
		if sts, ok := mId2BaseData[id]; ok {
			base := sts.PreClosePrice
			ra1 := getRa(hight, base)
			ra2 := getRa(low, base)
			diff := ra1 - ra2
			//fmt.Println("...", hidx, low, diff, base)
			if diff >= WarnRatio && diff < 20 {
				ch := "↑↑↑"
				if lidx == n-1 {
					ch = "↓↓↓"
				}
				mId2cnt[id] = time.Now().Unix()
				SendMsg2Listen(id, fmt.Sprintf("%s 在%d次交易中异动 %s\n%.2f%%  %.2f%%  %.2f%% \n", muban, tcnt, ch, getRa(low, base), getRa(hight, base), diff))
				fmt.Printf("%s 在%d次交易中异动 %s %.2f%%  %.2f%%  %.2f%%  %s %.2f %.2f %v\n", muban, tcnt, ch, getRa(low, base), getRa(hight, base), diff, id, low, hight, base)
				return false
			} else {
				fmt.Println("...", hidx, low, diff, base, hight, low)
			}
		}
	}
	return true
}

// 突破上限才比较
func checkUnActionMaxMin(r dataRes) {
	muban := ""
	if info, ok := mId2ConstInfo[r.Inst]; ok {
		muban = info.InstrumentName
	}
	if x, ok := mId2HL[r.Inst]; ok {
		for _, y := range r.QuoteData.DynaData {
			if y.HighestPrice > x.Max {
				if sts, ok := mId2BaseData[r.Inst]; ok {
					ratio1 := getRa(y.HighestPrice, sts.PreClosePrice)
					ratio2 := getRa(x.Max, sts.PreClosePrice)
					ratio3 := getRa(x.Min, sts.PreClosePrice)
					fmt.Println(ratio2/0.49, ratio1/0.49)
					if (ratio1/0.49 != ratio2/0.49) && (ratio1-ratio3 > WarnHL) {
						SendMsg2Listen(r.Inst, fmt.Sprintf("%s 新高↑↑↑\n%.2f%%  %.2f%%\n", muban, ratio1, ratio3))
						fmt.Printf("%s 新高↑↑↑\n%.2f%%  %.2f%% inst = %s %.2f %.2f\n", muban, ratio1, ratio3, r.Inst, x.Max, x.Min)
					}
				}
				x.Max = y.HighestPrice
			}
			if y.LowestPrice != 0 && y.LowestPrice < x.Min {
				if sts, ok := mId2BaseData[r.Inst]; ok {
					ratio1 := getRa(y.LowestPrice, sts.PreClosePrice)
					ratio2 := getRa(x.Min, sts.PreClosePrice)
					ratio3 := getRa(x.Max, sts.PreClosePrice)

					if (ratio1/0.49 != ratio2/0.49) && (ratio3-ratio1 > WarnHL) {
						SendMsg2Listen(r.Inst, fmt.Sprintf("%s 新低↓↓↓\n%.2f%%  %.2f%%\n", muban, ratio3, ratio1))
						fmt.Printf("%s 新低\n%.2f%%  %.2f%% inst = %s %.2f %.2f\n", muban, ratio1, ratio3, r.Inst, x.Max, x.Min)

					}
				}
				x.Min = y.LowestPrice
			}
		}
	} else {
		for _, x := range r.QuoteData.DynaData {
			if x.LowestPrice > 0 {

				mId2HL[r.Inst] = &HL{Max: x.HighestPrice, Min: x.LowestPrice}
				//	fmt.Println("--- ", r.Inst, x.HighestPrice, x.LowestPrice)
			} else {
				fmt.Println("err dyna:", r.Inst, x.HighestPrice, x.LowestPrice)
			}
		}
	}
}

func checkVra(id string, mul float64, run bool) bool {
	if !run {
		return false
	}

	return true
}
func GetList(id string) string {
	msg := "list:\n"
	for f, val := range getFollow(id).FollowsId {
		if val.Ok == false {
			continue
		}
		idx := len(mId2Tick[f]) - 1
		if idx < 0 {
			continue
		}
		if sts, ok := mId2BaseData[f]; ok {

			tick := mId2Tick[f][idx]
			ratio := getRa(tick.Price, sts.PreClosePrice)
			name := sts.InstrumentID
			if info, ok := mId2ConstInfo[f]; ok {
				name = info.InstrumentName
			}
			msg += fmt.Sprintf("%-10s  %.2f%%\n", name, ratio)
		}
	}
	return msg
}
func handleDyna(r dataRes) {
	for _, x := range r.QuoteData.DynaData {
		mId2Dyna[r.Inst] = append(mId2Dyna[r.Inst], x)
		checkUnActionMaxMin(r)
	}
}
func handleSTATISTICS(r dataRes) {
	for _, x := range r.QuoteData.StatisticsData {
		mId2BaseData[r.Inst] = x

	}
}
func handleStatic(r dataRes) {
	for _, x := range r.QuoteData.StaticData {
		mId2ConstInfo[r.Inst] = x
	}
}

func handleRes(r dataRes) {

	if r.ServiceType == "TICK" {
		handleTick(r)
	} else if r.ServiceType == "DYNA" {
		handleDyna(r)
	} else if r.ServiceType == "STATISTICS" {
		handleSTATISTICS(r)
	} else if r.ServiceType == "STATIC" {
		handleStatic(r)
	}

}
func RunWs() {
	wscon = make([]*websocket.Conn, wsc+1)
	for i := 0; i < wsc; i++ {
		go startws(i)
	}
	startws(wsc)
}
func startws(i int) {

	header := http.Header{
		"Accept-Language": []string{"zh-CN,zh;q=0.9"},
		"Accept-Encoding": []string{"gzip, deflate, br"},
		"Cache-Control":   []string{"no-cache"},
		"Host":            []string{"qas.sylapp.cn"},
		"Origin":          []string{"https://stock.9fzt.com"},
		"Pragma":          []string{"no-cache"},
		"User-Agent":      []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36"},
	}
	url := "wss://qas.sylapp.cn/quote"
	dl := websocket.Dialer{}
	conn, res, err := dl.Dial(url, header)
	//conn, res, err := websocket.DefaultDialer.Dial(url, header)
	fmt.Println("res ", res)
	if err != nil {
		fmt.Println(header)
		fmt.Println("建立链接失败！！", err)
		return
	} else {
		fmt.Println("connect success")
	}
	daya := Daya_json{
		OrgCode: "rh",
		Token:   "e9252a64-6ac8-4bf8-9725-6f106f682908",
		AppName: "htzpc",
		AppVer:  "v1.0.0",
		AppType: "pc",
		Tag:     "af91c60a-1acc-4150-9014-d11086cb9489",
	}

	err = conn.WriteJSON(daya)
	if err != nil {
		fmt.Println("json marshal error:", err)
	}
	go Ping(conn)
	wscon[i] = conn
	if err != nil {
		fmt.Println("json marshal error:", err)
	}
	if i == wsc {
		go ReLoad()
		DsMsg()
	}
	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		pong := Pong{}
		err = json.Unmarshal(b, &pong)
		if err == nil && pong.Code == "200" {
			fmt.Printf("%+v\n", pong)
			continue
		}
		r := dataRes{}
		err = json.Unmarshal(b, &r)
		if err != nil {
			fmt.Println("json.Unmarshal error ", err)
		} else {
			//fmt.Println("ws")
			handleRes(r)
		}
	}
}
