package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const OneHand = 100

type empty struct {
}

var reqid int32 = 1
var mId2Ratio map[string]float64

// 开盘时间
var mId2BaseData map[string]StatisticType

// 基本信息
var mId2ConstInfo map[string]StaticType

// 监听者
var mId2Listener map[string]map[string]empty
var cconn *websocket.Conn

func PostSTATISTICS(id string) {
	addr := "sh"
	dataJson := Data_json{
		SubType:     "SUBON",
		Inst:        id,
		Market:      addr,
		ServiceType: "STATISTICS",
		ReqID:       int(atomic.AddInt32(&reqid, 1)),
	}
	cconn.WriteJSON(dataJson)
	addr = "sz"
	cconn.WriteJSON(dataJson)
}
func PostTick(id string) {
	addr := "sh"
	dataJson := Data_json{
		SubType:     "SUBON",
		Inst:        id,
		Market:      addr,
		ServiceType: "TICK",
		ReqID:       int(atomic.AddInt32(&reqid, 1)),
	}
	cconn.WriteJSON(dataJson)
	addr = "sz"
	cconn.WriteJSON(dataJson)
}
func PostStatic(id string) {
	addr := "sh"
	dataJson := Data_json{
		SubType:     "SUBON",
		Inst:        id,
		Market:      addr,
		ServiceType: "STATIC",
		ReqID:       int(atomic.AddInt32(&reqid, 1)),
	}
	cconn.WriteJSON(dataJson)
	addr = "sz"
	cconn.WriteJSON(dataJson)
}
func PostDyna(id string) {
	addr := "sh"
	dataJson := Data_json{
		SubType:     "SUBON",
		Inst:        id,
		Market:      addr,
		ServiceType: "DYNA",
		ReqID:       int(atomic.AddInt32(&reqid, 1)),
	}
	cconn.WriteJSON(dataJson)
	addr = "sz"
	cconn.WriteJSON(dataJson)

}

func Post(name, id string) {
	//当前没有post过
	if _, ok := mId2ConstInfo[id]; !ok {
		PostStatic(id)
		PostDyna(id)
		PostSTATISTICS(id)
		PostTick(id)
	} else {
		mId2Listener[id][name] = empty{}
	}
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
	if info, ok := mId2ConstInfo[r.Inst]; ok {
		muban = info.InstrumentName + "\n"
	}
	b := false
	for _, x := range r.QuoteData.TickData {
		if x.Volume > 200*OneHand {
			base := x.Price
			if sts, ok := mId2BaseData[r.Inst]; ok {
				base = sts.PreClosePrice
			}
			b = true
			v = append(v, x.Volume)
			str = append(str, fmt.Sprintf("%s   %g   %d  %.2f%%\n", r.Inst, x.Price, x.Volume/OneHand, getRa(x.Price, base)))
		}
	}
	if b {
		n := len(v)
		for name, _ := range mId2Listener[r.Inst] {
			needlen := getFoller(name).Id[r.Inst]
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
	}
}
func handleDyna(r dataRes) {
	for _, x := range r.QuoteData.DynaData {
		if _, ok := mId2Ratio[r.Inst]; ok {
			//fmt.Println(v)
		} else {
			mId2Ratio[r.Inst] = x.LastPrice
		}
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

	} else if r.ServiceType == "STATISTICS" {
		handleSTATISTICS(r)
	} else if r.ServiceType == "STATIC" {
		handleStatic(r)
	}

}
func InitWs() {
	mId2Ratio = map[string]float64{}
	mId2BaseData = map[string]StatisticType{}
	mId2ConstInfo = map[string]StaticType{}
	mId2Listener = map[string]map[string]empty{}
	mNameFollor = map[string]*Follow{}
}
func Run() {

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
	InitWs()
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
	cconn = conn
	go Ping(conn)
	data := Data_json{
		SubType:     "SUBON",
		Inst:        "603912",
		Market:      "sh",
		ServiceType: "TICK",
		ReqID:       0,
	}

	err = conn.WriteJSON(data)
	if err != nil {
		fmt.Println("json marshal error:", err)
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
			fmt.Printf("ws: %+v\n", r)
			handleRes(r)
		}
	}
}
