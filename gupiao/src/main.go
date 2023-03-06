package main

import (
	"net/http"
	"time"
)

var StartWs chan struct{}

func InitVal() {
	mId2Dyna = map[string][]DynaType{}
	mId2BaseData = map[string]StatisticType{}
	mId2ConstInfo = map[string]StaticType{}
	mId2Listener = map[string]map[string]empty{}
	mIdFollow = map[string]*Follow{}
	mId2HL = map[string]*HL{}
	mId2Tick = map[string][]TickType{}
	mId2LB = map[string]*VRa{}
	MsgChan = make(chan MsgType, 100)
	StartWs = make(chan struct{}, 1)
	client = &http.Client{Timeout: 5 * time.Second}
	ddMsg = make(chan string, 100)
	mId2Post = map[string]empty{}
	mId2cnt = map[string]int64{}
	mId2Time = map[string]int64{}
}

func main() {

	InitVal()
	InitRedis()
	go RunWechat()
	go RecvDDMsg()
	<-StartWs

	RunWs()
	//s := "a  000  sz"
	//fmt.Println(stringSplit(s, ' '))
}
