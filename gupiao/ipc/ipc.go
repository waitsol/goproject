package ipc

import "main/proto/go/main/pb3"

type MsgType struct {
	Id  string
	Msg string
}

var MsgChan chan MsgType
var ReqChan chan pb3.PacketReq
var RspChan chan pb3.PacketRsp

func init() {
	MsgChan = make(chan MsgType, 100)
	ReqChan = make(chan pb3.PacketReq, 100)
	RspChan = make(chan pb3.PacketRsp, 100)
}
