package user

import "sync"

// 监听者
var SyncId2Listener sync.Map

type FollowInfo struct {
	Num   int
	MaxRa float64
	MinRa float64
}

func DelFollow(gid, uid string) {

	load, ok := SyncId2Listener.Load(gid)
	if ok {
		listens := load.(map[string]*FollowInfo)
		delete(listens, uid)
	}
}
func AddFollow(gid, uid string, val int) {

	load, ok := SyncId2Listener.LoadOrStore(gid, map[string]*FollowInfo{uid: &FollowInfo{Num: val, MaxRa: 200, MinRa: -200}})
	if ok == true {
		x := load.(map[string]*FollowInfo)
		if x[uid] == nil {
			x[uid] = &FollowInfo{}
		}
		x[uid].Num = val
	}
}

func SetFollowMaxRa(gid, uid string, val float64) {
	load, ok := SyncId2Listener.LoadOrStore(gid, map[string]*FollowInfo{uid: &FollowInfo{Num: 2000, MaxRa: val}})
	if ok == true {
		x := load.(map[string]*FollowInfo)
		x[uid].MaxRa = val
	}

}
func SetFollowMinRa(gid, uid string, val float64) {
	load, ok := SyncId2Listener.LoadOrStore(gid, map[string]*FollowInfo{uid: &FollowInfo{Num: 2000, MaxRa: val}})
	if ok == true {
		x := load.(map[string]*FollowInfo)
		x[uid].MinRa = val
	}
}
