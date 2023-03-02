package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync/atomic"
	"time"
)

const OneHand = 100

var reqid int32 = 1
var mName2Ratio map[string]float64
var cconn *websocket.Conn

type Data_json struct {
	SubType     string `json:"SubType"`
	ReqID       int    `json:"ReqID"`
	Inst        string `json:"Inst"`
	Market      string `json:"Market"`
	ServiceType string `json:"ServiceType"`
}

type Daya_json struct {
	OrgCode string `json:"OrgCode"`
	Token   string `json:"Token"`
	AppName string `json:"AppName"`
	AppVer  string `json:"AppVer"`
	AppType string `json:"AppType"`
	Tag     string `json:"Tag"`
}
type DynaType struct {
	TradingDay     int       `json:"TradingDay"`
	Time           int       `json:"Time"`
	HighestPrice   float64   `json:"HighestPrice"`
	LowestPrice    float64   `json:"LowestPrice"`
	LastPrice      float64   `json:"LastPrice"`
	Volume         int       `json:"Volume"`
	Amount         float64   `json:"Amount"`
	TickCount      int       `json:"TickCount"`
	BuyPrice       []float64 `json:"BuyPrice"`
	BuyVolume      []int     `json:"BuyVolume"`
	SellPrice      []float64 `json:"SellPrice"`
	SellVolume     []int     `json:"SellVolume"`
	AveragePrice   float64   `json:"AveragePrice"`
	Wk52High       float64   `json:"Wk52High"`
	Wk52Low        float64   `json:"Wk52Low"`
	PERatio        float64   `json:"PERatio"`
	OrderDirection int       `json:"OrderDirection"`
	BidPrice       float64   `json:"BidPrice"`
	AskPrice       float64   `json:"AskPrice"`
	TurnoverRate   float64   `json:"TurnoverRate"`
	SA             float64   `json:"SA"`
	LimitUp        float64   `json:"LimitUp"`
	LimitDown      float64   `json:"LimitDown"`
	CirStock       float64   `json:"CirStock"`
	TotStock       float64   `json:"TotStock"`
	CirVal         float64   `json:"CirVal"`
	TotVal         float64   `json:"TotVal"`
	NAV            float64   `json:"NAV"`
	Ratio          float64   `json:"Ratio"`
	Committee      float64   `json:"Committee"`
	PES            float64   `json:"PES"`
	WP             int       `json:"WP"`
	NP             int       `json:"NP"`
	LastTradeVol   int       `json:"LastTradeVol"`
	YearUpDown     float64   `json:"YearUpDown"`
	KindsUpdown    struct {
		FiveMinsUpdown  float64 `json:"FiveMinsUpdown"`
		ThreeMinsUpdown float64 `json:"ThreeMinsUpdown"`
		OneMinsUpdown   int     `json:"OneMinsUpdown"`
		MinUpdown2      int     `json:"MinUpdown2"`
		MinUpdown4      int     `json:"MinUpdown4"`
	} `json:"KindsUpdown"`
	Updown               float64 `json:"Updown"`
	NextDayPreClosePrice float64 `json:"NextDayPreClosePrice"`
	ExchangeID           string  `json:"ExchangeID"`
	InstrumentID         string  `json:"InstrumentID"`
	TTM                  float64 `json:"TTM"`
}
type TickType struct {
	TradingDay int     `json:"TradingDay"`
	ID         int     `json:"ID"`
	Time       int     `json:"Time"`
	Price      float64 `json:"Price"`
	Volume     int     `json:"Volume"`
	Property   int     `json:"Property"`
	Virtual    int     `json:"virtual"`
}
type InstStatusType struct {
	StatusType   int    `json:"StatusType"`
	ExchangeID   string `json:"ExchangeID"`
	InstrumentID string `json:"InstrumentID"`
}
type KlineType struct {
	TradingDay int     `json:"TradingDay"`
	Time       int     `json:"Time"`
	High       float64 `json:"High"`
	Open       float64 `json:"Open"`
	Low        float64 `json:"Low"`
	Close      float64 `json:"Close"`
	Volume     int     `json:"Volume"`
	Amount     float64 `json:"Amount"`
	TickCount  int     `json:"TickCount"`
}
type StatisticType struct {
	TradingDay int `json:"TradingDay"`
	//昨天价格
	PreClosePrice float64 `json:"PreClosePrice"`
	//开盘价
	OpenPrice       float64 `json:"OpenPrice"`
	UpperLimitPrice float64 `json:"UpperLimitPrice"`
	LowerLimitPrice float64 `json:"LowerLimitPrice"`
	ExchangeID      string  `json:"ExchangeID"`
	InstrumentID    string  `json:"InstrumentID"`
}

type MinType struct {
	TradingDay       int     `json:"TradingDay"`
	Time             int     `json:"Time"`
	Price            float64 `json:"Price"`
	Volume           int     `json:"Volume"`
	UnmismatchVolume int     `json:"UnmismatchVolume,omitempty"`
	UnmismatchFlag   int     `json:"UnmismatchFlag,omitempty"`
}
type dataRes struct {
	Market      string `json:"Market"`
	Inst        string `json:"Inst"`
	ServiceType string `json:"ServiceType"`
	SubType     string `json:"SubType"`
	ReqID       int    `json:"ReqID"`
	QuoteData   struct {
		//动态数据
		DynaData []DynaType `json:"DynaData"`
		//交易买单
		TickData []TickType `json:"TickData"`
		//
		InstStatusData []InstStatusType `json:"InstStatusData"`
		KlineData      []KlineType      `json:"KlineData"`
		StatisticsData []StatisticType  `json:"StatisticsData"`
		//早盘
		MinData []MinType `json:"MinData"`
	} `json:"QuoteData"`
}
type PingType struct {
	ServiceType string `json:"ServiceType"`
}
type Pong struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}
type T struct {
	Market      string `json:"Market"`
	Inst        string `json:"Inst"`
	ServiceType string `json:"ServiceType"`
	SubType     string `json:"SubType"`
	ReqID       int    `json:"ReqID"`
	QuoteData   struct {
		TickData []struct {
			TradingDay int     `json:"TradingDay"`
			ID         int     `json:"ID"`
			Time       int     `json:"Time"`
			Price      float64 `json:"Price"`
			Volume     int     `json:"Volume"`
			Property   int     `json:"Property"`
			Virtual    int     `json:"virtual"`
		} `json:"TickData"`
	} `json:"QuoteData"`
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
func handleTick(r dataRes) {
	for _, x := range r.QuoteData.TickData {
		if x.Volume > 200*OneHand {
			SendMsg("wm", fmt.Sprintf("%s   %g   %d\n", r.Inst, x.Price, x.Volume/OneHand))
		}
	}
}
func handleDyna(r dataRes) {
	for _, x := range r.QuoteData.DynaData {
		if _, ok := mName2Ratio[r.Inst]; ok {
			//fmt.Println(v)
		} else {
			mName2Ratio[r.Inst] = x.LastPrice
		}
	}
}
func Post(dataJson Data_json) {
	fmt.Println("post ", dataJson)
	dataJson.ReqID = int(atomic.AddInt32(&reqid, 1))
	cconn.WriteJSON(dataJson)
}
func handleRes(r dataRes) {

	if r.ServiceType == "TICK" {
		handleTick(r)
	} else if r.ServiceType == "DYNA" {

	}

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
	mName2Ratio = map[string]float64{}
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
