package redis

import (
	"encoding/json"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"main/dingding"
	"os"
)

var cliRedis *redis.Client

const HKey = "UserInfo"
const AllID = "allgpid"

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
		panic("redis no url")
	}
	dingding.DDURL = url
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
