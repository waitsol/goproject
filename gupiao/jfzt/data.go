package jfzt

import "github.com/gorilla/websocket"

type HL struct {
	Min float64
	Max float64
}

var mId2TurnoverRate = map[string][]float64{}

type KvInfo struct {
	val   int64
	cnt   int64
	begin float64
}

func (this *KvInfo) Reset() {
	this.cnt = 0
	this.val = 0
}

// 交易情况
type TsInfo struct {
	vec []KvInfo
}

type Stock struct {
	// 最高价和最低价
	hl HL
	// 开盘时间
	baseData StatisticType
	// 基本信息
	constInfo StaticType
	// 最近买单
	tick []TickType
	// 最后价格
	dyna []DynaType
	// 量比
	lb    *VRa
	ttime int64
	cnt   int64

	//上涨情况
	up KvInfo
	//下跌情况
	down KvInfo

	//交易信息
	buy  TsInfo
	sell TsInfo
	//买卖flag
	msFlag  []bool
	bigSell int64
	bigBuy  int64
}

// 每个ws一个 ws收发单协程 协程安全
type WsSet struct {
	start bool
	// 最高价和最低价
	mId2HL map[string]*HL
	// 开盘时间
	mId2BaseData map[string]StatisticType
	// 基本信息
	mId2ConstInfo map[string]StaticType
	// 最近买单
	mId2Tick map[string][]TickType
	// 最后价格
	mId2Dyna map[string][]DynaType

	mId2LB   map[string]*VRa
	mId2Time map[string]int64
	mId2cnt  map[string]int64
	stock    map[string]*Stock
	conn     *websocket.Conn
}

func (this *WsSet) Init() {
	this.mId2Dyna = map[string][]DynaType{}
	this.mId2BaseData = map[string]StatisticType{}
	this.mId2ConstInfo = map[string]StaticType{}
	this.mId2HL = map[string]*HL{}
	this.mId2Tick = map[string][]TickType{}
	this.mId2LB = map[string]*VRa{}
	this.mId2cnt = map[string]int64{}
	this.mId2Time = map[string]int64{}
	this.stock = map[string]*Stock{}
	this.start = true

}

func (this *WsSet) Reset() {
	this.mId2Dyna = map[string][]DynaType{}
	this.mId2HL = map[string]*HL{}
	this.mId2Tick = map[string][]TickType{}
	this.mId2LB = map[string]*VRa{}
	this.mId2cnt = map[string]int64{}
	this.mId2Time = map[string]int64{}
	this.stock = map[string]*Stock{}
	this.start = true

}
func (this *WsSet) Stop() {

	this.start = false

}
