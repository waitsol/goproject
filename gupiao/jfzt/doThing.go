package jfzt

func DelFollow(gid, uid string) {

	load, ok := SyncId2Listener.Load(gid)
	if ok {
		listens := load.(map[string]*followInfo)
		delete(listens, uid)
	}
}
func AddFollow(gid, uid string, val int) {

	load, ok := SyncId2Listener.LoadOrStore(gid, map[string]*followInfo{uid: &followInfo{num: val, maxRa: 200, minRa: -200}})
	if ok == true {
		x := load.(map[string]*followInfo)
		if x[uid] == nil {
			x[uid] = &followInfo{}
		}
		x[uid].num = val
	}
}

func SetFollowMaxRa(gid, uid string, val float64) {
	load, ok := SyncId2Listener.LoadOrStore(gid, map[string]*followInfo{uid: &followInfo{num: 2000, maxRa: val}})
	if ok == true {
		x := load.(map[string]*followInfo)
		x[uid].maxRa = val
	}

}
func SetFollowMinRa(gid, uid string, val float64) {
	load, ok := SyncId2Listener.LoadOrStore(gid, map[string]*followInfo{uid: &followInfo{num: 2000, maxRa: val}})
	if ok == true {
		x := load.(map[string]*followInfo)
		x[uid].minRa = val
	}
}

func SendMsg(id, msg string, send_group bool) {
	//MsgChan <- MsgType{id, msg}
	//dingding.DdMsg <- dingding.DDMsgType{Id: id, Msg: msg}
	m := map[string]interface{}{}
	m["id"] = id
	if send_group {
		m["group_id"] = "853312133"
	}
	im.SendMsg(msg, m)
}
