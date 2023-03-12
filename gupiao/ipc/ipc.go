package ipc

type MsgType struct {
	Id  string
	Msg string
}

var MsgChan chan MsgType

func init() {
	MsgChan = make(chan MsgType, 100)
}
