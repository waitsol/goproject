package redis

import (
	"encoding/json"
	"fmt"
	"main/dingding"
	"os"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var cliRedis *redis.Client

const HKey = "UserInfo"
const AllID = "allgpid"
const name2id = "name2id"
const id2name = "id2name"

type RedisCfg struct {
	Addr string `json:"addr"`
	Mima string `json:"mima"`
}

func init() {
	f, err := os.Open("redis.json")
	if err != nil {
		log.Error(err)
		log.Exit(-1)
	}
	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil {
		log.Error(err)
		log.Exit(-1)
	}
	r := RedisCfg{}
	err = json.Unmarshal(buf[:n], &r)
	if err != nil {
		log.Error(err)
		log.Exit(-1)
	}
	cliRedis = redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Mima, // no password set
		DB:       0,      // use default DB
	})
	if nil != cliRedis.Ping().Err() {
		panic("redis connect error")
	}
	url, err := cliRedis.Get("wmurl").Result()
	if err != nil || len(url) == 0 {
		//panic("redis no url")
	}
	dingding.DDURL = url
	dingding.DDURL = "https://oapi.dingtalk.com/robot/send?access_token=6982ee152bc9709ce56301cb8fc3de7d2b63a31c73e19dbca3b6f360374f85ad"
	log.Info(dingding.DDURL)
}
func SaveUser(name string, bitdata []byte) {

	if nil != cliRedis.HSet(HKey, name, bitdata).Err() {
		log.Error("save failed ", name)
	}

}
func LoadUser(id string) (string, error) {
	data, err := cliRedis.HGet(HKey, id).Result()
	return data, err
}

func LoadFollow() (map[string]string, error) {
	return cliRedis.HGetAll(HKey).Result()
}
func GetInfoFromRedis(gid string) (string, error) {
	return cliRedis.HGet(AllID, gid).Result()
}
func LoadAll() (map[string]string, error) {
	return cliRedis.HGetAll(AllID).Result()
}
func GetDQ(gid string) string {
	dq, err := cliRedis.HGet(AllID, gid).Result()
	if err == nil {
		return dq
	}
	return "sh"
}
func SaveTurnoverRate(gid string, ra float64) {
	cliRedis.RPush(fmt.Sprintf("TurnoverRate.%v", gid), ra)
	if cliRedis.LLen(fmt.Sprintf("TurnoverRate.%v", gid)).Val() > 40 {
		cliRedis.LPop(fmt.Sprintf("TurnoverRate.%v", gid))
	}
	return
}
func LoadTurnoverRate(gid string) []string {
	x, _ := cliRedis.LRange(fmt.Sprintf("TurnoverRate.%v", gid), 0, -1).Result()
	return x
}
func FixData(gid string) {
	if cliRedis.LLen(fmt.Sprintf("TurnoverRate.%v", gid)).Val() == 5 {
		data_, _ := cliRedis.RPop(fmt.Sprintf("TurnoverRate.%v", gid)).Result()
		cliRedis.RPop(fmt.Sprintf("TurnoverRate.%v", gid)).Result()
		cliRedis.RPush(fmt.Sprintf("TurnoverRate.%v", gid), data_)
	}
	return
}
func Name2Id(name string) string {
	xx, err := cliRedis.HGet(name2id, name).Result()
	if err != nil {
		return ""
	}
	return xx
}
func Id2Name(id string) string {
	xx, err := cliRedis.HGet(id2name, id).Result()
	if err != nil {
		return ""
	}
	return xx
}
