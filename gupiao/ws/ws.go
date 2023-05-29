package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/waitsol/golib"
	"main/redis"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	OneHand   = 100
	WarnHL    = 1
	WarnRatio = 230
	LB        = 300
	WarnCheck = 2
)
const WSC = 4

var checkSecs = []int64{3, 10, 30, 60, 120, 300}
var checkCnts = []int{10, 20, 30, 50, 100}

var mId2Post map[string]Empty

func init() {
	mId2Post = map[string]Empty{}
}

type Empty struct {
}

var reqid int32 = 1

// 监听者
var SyncId2Listener sync.Map

var MGR []*WsSet

func hx(id string) int {
	x, _ := strconv.Atoi(id)
	return x % (WSC)
}

func GetRa(cur, base float64) float64 {
	return (cur - base) / base * 100
}
func clock(f func()) {
	t1 := time.Now().UnixMilli()
	f()
	t2 := time.Now().UnixMilli()
	fmt.Println(t2 - t1)
}

func RunWs() {

	x, _ := redis.LoadAll()
	for k, _ := range x {
		data := redis.LoadTurnoverRate(k)
		for _, ra := range data {
			f, _ := strconv.ParseFloat(ra, 64)
			mId2TurnoverRate[k] = append(mId2TurnoverRate[k], f)
		}
	}
	MGR = make([]*WsSet, WSC)
	for i := 0; i < WSC; i++ {
		MGR[i] = &WsSet{}
		golib.Go(func() {
			startws(i)
		})
		time.Sleep(1 * time.Second)
	}
	//定时消息
	DsMsg()
}
func load(wsidx int) {
	//加载所有股票
	data, err := redis.LoadAll()
	if err == nil {
		for gpid, dq := range data {
			if hx(gpid) == wsidx {
				GetMgr(gpid).Post(gpid, dq)
			}
		}
	}

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
	log.Info("res ", res)
	if err != nil {
		log.Error(header)
		log.Error("建立链接失败！！", err)
		return
	} else {
		log.Info("connect success")
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
		log.Error("json marshal error:", err)
	}
	MGR[i].Init()
	MGR[i].conn = conn
	if time.Now().Hour() >= 24 {
		MGR[i].Stop()
	}

	stopc := make(chan bool)
	golib.Go(func() {
		MGR[i].Ping(stopc)
	})
	load(i)
	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			log.Error(err)
			//ws断开重启
			stopc <- false
			startws(i)
			return
		}
		pong := Pong{}
		err = json.Unmarshal(b, &pong)
		if err == nil && pong.Code == "200" {
			log.Infof("%+v\n", pong)
			continue
		}
		r := dataRes{}
		err = json.Unmarshal(b, &r)
		if err != nil {
			log.Error("json.Unmarshal error ", err)
		} else {
			//fmt.Printf("ws: %+v\n", r)
			MGR[i].handleRes(r)
		}
	}
}
