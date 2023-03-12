package wechat

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"main/redis"
)

func SaveUserFollow(follow Follow) {
	bitdata, err := json.Marshal(follow)
	if err == nil {
		redis.SaveUser(follow.Id, bitdata)
	} else {
		log.Error("save failed ", follow.Id)
	}
}
func GetUserFollowFromRedis(id string) *Follow {
	data, err := redis.LoadUser(id)
	x := &Follow{}
	x.FollowsId = map[string]*FollowSt{}
	if err == nil && json.Unmarshal([]byte(data), x) != nil {
		return x
	}
	x.Id = id
	return x
}
func checkGpNum(s string) bool {

	_, err := redis.GetInfoFromRedis(s)
	return err == nil
}
