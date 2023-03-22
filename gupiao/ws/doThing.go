package ws

import (
	"main/dingding"
)

func DelFollow(gid, uid string) {

	load, ok := SyncId2Listener.Load(gid)
	if ok {
		listens := load.(map[string]int)
		delete(listens, uid)
	}
}
func AddFollow(gid, uid string, val int) {

	load, ok := SyncId2Listener.LoadOrStore(gid, map[string]int{uid: val})
	if ok == true {
		x := load.(map[string]int)
		x[uid] = val
	}
}
func SendMsg(id, msg string) {
	//MsgChan <- MsgType{id, msg}
	dingding.DdMsg <- dingding.DDMsgType{Id: id, Msg: msg}
}
