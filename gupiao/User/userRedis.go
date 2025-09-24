package user

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
func getRealGpNum(s string) string {

	num := redis.Name2Id(s)
	return num
}
func Id2Name(s string) string {
	return redis.Id2Name(s)
}
