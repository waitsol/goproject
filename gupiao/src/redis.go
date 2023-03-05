package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
)

var cliRedis *redis.Client

const HKey = "UserInfo"
const AllID = "allgpid"

func InitRedis() {
	cliRedis = redis.NewClient(&redis.Options{
		Addr:     "218.78.68.173:1124",
		Password: "042199ww", // no password set
		DB:       0,  // use default DB
	})
	if nil != cliRedis.Ping().Err() {
		panic("redis connect error")
	}
	url, err := cliRedis.Get("wmurl").Result()
	if err != nil || len(url) == 0 {
		panic("redis no url")
	}
	DDURL = url
	fmt.Println(DDURL)
}
func SaveUserFollow(name string, follow Follow) {
	bitdata, err := json.Marshal(follow)
	if err == nil {
		if nil != cliRedis.HSet(HKey, name, bitdata).Err() {
			fmt.Println("save failed ", name)
		}
	} else {
		fmt.Println("save failed ", name)
	}
}
func GetUserFollow(id string) *Follow {
	data, err := cliRedis.HGet(HKey, id).Result()
	x := &Follow{}
	if err == nil && json.Unmarshal([]byte(data), x) != nil {
		return x
	}
	x.Id = id
	x.FollowsId = map[string]*FollowSt{}
	return x
}
func ClearFollowById(id string) {
	cliRedis.HDel(HKey, id)
}

func ReLoad() {
	data, err := cliRedis.HGetAll(HKey).Result()
	if err == nil {
		for wechatid, userinfo := range data {
			f := &Follow{}
			if json.Unmarshal([]byte(userinfo), f) == nil {
				mIdFollow[wechatid] = f
				for fid, _ := range f.FollowsId {
					Post(wechatid, fid)
				}
			}
		}
	}
	LoadAll()
}

func LoadAll() {
	data, err := cliRedis.HGetAll(AllID).Result()
	if err == nil {
		for gpid, _ := range data {
			PostStatic(gpid)
			PostSTATISTICS(gpid)
			PostTick(gpid)
			PostDyna(gpid)
		}
	}

}
