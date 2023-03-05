package main

import (
	"strconv"
)

var mIdFollow map[string]*Follow

type FollowSt struct {
	WarnMsg int  `json:"WarnMsg"`
	Ok      bool `json:"Ok"`
}
type Follow struct {
	FollowsId map[string]*FollowSt `json:"FollowsId"` //关注的股票id
	Id        string               `json:"Id"`
}

func (this *Follow) follow(id string) {
	//如果关注了
	if x, ok := this.FollowsId[id]; ok && x.Ok {
		//退出监听
		delete(mId2Listener[id], this.Id)
		delete(mIdFollow, this.Id)
		this.FollowsId[id].Ok = false
	} else {
		if _, ok := this.FollowsId[id]; !ok {
			this.FollowsId[id] = &FollowSt{Ok: true, WarnMsg: 400}
		}
		this.FollowsId[id].Ok = true

		Post(this.Id, id)
	}

	SaveUserFollow(this.Id, *this)
}

func (this *Follow) HandleMessage(msg string) (bool, string) {
	v := stringSplit(msg, ' ')
	if len(v) == 0 {
		return false, ""
	}
	if checkGpNum(v[0]) {
		this.follow(v[0])
		return true, "苟富贵"
	} else if v[0] == "set" {
		if len(v) == 3 && checkGpNum(v[1]) {
			x, err := strconv.Atoi(v[2])
			if err != nil || x < 200 {
				return true, "err " + v[2]
			}
			if this.FollowsId[v[1]] == nil {
				this.FollowsId[v[1]] = &FollowSt{Ok: false}
			}
			this.FollowsId[v[1]].WarnMsg = x * OneHand
			return true, "ok"
		} else {
			return true, "err args"
		}
	} else if v[0] == "list" {
		_sendMsg(this.Id, GetList(this.Id))
	} else if v[0] == "clear" {
		ClearFollowById(this.Id)
	}
	return false, ""
}
func getFollow(id string) *Follow {
	if mIdFollow == nil {
		mIdFollow = map[string]*Follow{}
	}
	if x, ok := mIdFollow[id]; ok {
		return x
	}
	x := GetUserFollow(id)
	mIdFollow[id] = x
	return x
}

// 分割文本字符串
func stringSplit(text string, ic byte) []string {
	vs := []string{}

	b := -1
	for i, c := range text {
		if c == rune(ic) {
			if b != -1 {
				vs = append(vs, string(text[b:i]))
				b = -1
			}
		} else if b == -1 {
			b = i
		}
	}
	if b != -1 {
		vs = append(vs, string(text[b:]))
	}
	return vs
}
func checkGpNum(s string) bool {
	if len(s) != 6 {
		return false
	}
	for _, c := range s {
		if c > '9' || c < '0' {
			return false
		}
	}
	return true
}
