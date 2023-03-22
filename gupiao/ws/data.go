package ws

import "github.com/gorilla/websocket"

type HL struct {
	Min float64
	Max float64
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
	// 量比
	mId2LB   map[string]*VRa
	mId2Time map[string]int64
	mId2cnt  map[string]int64
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
	this.start = true

}

func (this *WsSet) Reset() {
	this.mId2Dyna = map[string][]DynaType{}
	this.mId2HL = map[string]*HL{}
	this.mId2Tick = map[string][]TickType{}
	this.mId2LB = map[string]*VRa{}
	this.mId2cnt = map[string]int64{}
	this.mId2Time = map[string]int64{}
	this.start = true

}
func (this *WsSet) Stop() {

	this.start = false

}
