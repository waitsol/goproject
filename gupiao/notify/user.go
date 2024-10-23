package notify

import (
	"encoding/json"
	"fmt"
	"main/jfzt"
	"main/redis"
	"math"
	"runtime/debug"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// 每个人关注那些 key roleid
var MIdFollow map[string]*Follow

func init() {
	MIdFollow = map[string]*Follow{}
}

type FollowSt struct {
	WarnMsg int  `json:"WarnMsg"`
	Ok      bool `json:"Ok"`
}
type Follow struct {
	FollowsId map[string]*FollowSt `json:"FollowsId"` //关注的股票id
	Id        string               `json:"Id"`
}

func (this *Follow) follow(id string) string {
	//如果关注了
	if x, ok := this.FollowsId[id]; ok && x.Ok {
		//退出WS监听
		jfzt.DelFollow(id, this.Id)
		//关闭自己列表监听
		this.FollowsId[id].Ok = false
		//保存自己信息
		SaveUserFollow(*this)
		return "取关成功"
	} else {
		if _, ok := this.FollowsId[id]; !ok {
			this.FollowsId[id] = &FollowSt{Ok: true, WarnMsg: 2000}
		}
		this.FollowsId[id].Ok = true
		if this.FollowsId[id].WarnMsg < 1000 {
			this.FollowsId[id].WarnMsg = 1000
		}

		jfzt.PostById(this.Id, id, this.FollowsId[id].WarnMsg)
	}

	SaveUserFollow(*this)
	return "关注成功"
}

func (this *Follow) HandleMessage(msg string) (bool, string) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			debug.PrintStack()
		}
	}()
	v := stringSplit(msg, ' ')
	if len(v) == 0 {
		return false, "填的什么玩意"
	}
	if v[0] == "/u" && len(v) == 2 && getRealGpNum(v[1]) != "" {

		return true, this.follow(getRealGpNum(v[1]))
	} else if v[0] == "/set" {
		num := getRealGpNum(v[1])
		if len(v) == 3 && num != "" {
			x, err := strconv.Atoi(v[2])
			if err != nil || x < 200 {
				return true, "err " + v[2]
			}
			if this.FollowsId[num] == nil {
				this.FollowsId[num] = &FollowSt{Ok: false}
			}
			this.FollowsId[num].WarnMsg = x
			SaveUserFollow(*this)
			jfzt.AddFollow(num, this.Id, this.FollowsId[num].WarnMsg)
			return true, "ok"
		} else {
			return true, "err args"
		}
	} else if v[0] == "/list" {
		return true, this.getList()
	} else if v[0] == "/clear" {
		this.clearUserFollow(this.Id)
		return true, "ok"
	} else if v[0] == "/Info" {
		result := "\n"
		for k, v := range this.FollowsId {
			result = fmt.Sprintf("%s%s : %d 监听:%v\n", result, Id2Name(k), v.WarnMsg, v.Ok)
		}
		return true, result
	} else if v[0] == "/min" {
		if len(v) == 3 && getRealGpNum(v[1]) != "" {
			x, err := strconv.Atoi(v[2])
			if err != nil || math.Abs(float64(x)) > 21 {
				return true, "err " + v[2]
			}

			jfzt.SetFollowMinRa(getRealGpNum(v[1]), this.Id, float64(x))
			return true, "ok"
		} else {
			return true, "err args"
		}
	} else if v[0] == "/max" {
		if len(v) == 3 && getRealGpNum(v[1]) != "" {
			x, err := strconv.Atoi(v[2])
			if err != nil || math.Abs(float64(x)) > 21 {
				return true, "err " + v[2]
			}
			jfzt.SetFollowMaxRa(getRealGpNum(v[1]), this.Id, float64(x))
			return true, "ok"
		} else {
			return true, "err args"
		}
	}
	return false, "填的什么玩意"
}
func GetFollow(id string) *Follow {
	if MIdFollow == nil {
		MIdFollow = map[string]*Follow{}
	}
	if x, ok := MIdFollow[id]; ok {
		return x
	}
	x := GetUserFollowFromRedis(id)
	MIdFollow[id] = x
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

func (this *Follow) clearUserFollow(uid string) {
	//先删除自己关注的
	for gid, v := range this.FollowsId {
		v.Ok = false
		jfzt.DelFollow(gid, uid)
	}
	SaveUserFollow(*this)
}
func (this *Follow) getList() string {
	v := []string{}
	for id, x := range this.FollowsId {
		if x.Ok {
			v = append(v, id)
		}
	}
	return jfzt.GetList(v)
}

// 处理关注列表
func loadFollow() {
	data, err := redis.LoadFollow()
	if err == nil {
		for wechatid, userinfo := range data {
			f := &Follow{}
			if json.Unmarshal([]byte(userinfo), f) == nil {
				MIdFollow[wechatid] = f
				for gid, v := range f.FollowsId {
					if v.Ok {
						jfzt.AddFollow(gid, wechatid, v.WarnMsg)
					}
				}
			}
		}
	}
}
