package main

import (
	"main/TcpNet"
	user "main/User"
	"main/jfzt"
	_ "main/log"
	"main/notify"
	_ "main/onebot11"
	_ "main/redis"

	log "github.com/sirupsen/logrus"
	"github.com/waitsol/golib"
)

var StartWs chan struct{}

func InitVal() {

	StartWs = make(chan struct{}, 1)

}
func simInfo(code string, all, base, fund int) {
	user.GetDefaultFollow().FollowsId[code] = &user.FollowSt{}
	user.GetDefaultFollow().FollowsId[code].All = all
	user.GetDefaultFollow().FollowsId[code].BaseHave = base
	user.GetDefaultFollow().FollowsId[code].Fund = fund
}
func main() {

	golib.Go(TcpNet.Run)
	user.GetDefaultFollow().AllFund = 100000
	user.GetDefaultFollow().FollowsId = make(map[string]*user.FollowSt)
	simInfo("601606", 1, 8, 30000)
	simInfo("600201", 63, 63, 10000)
	simInfo("600202", 8, 8, 10000)
	simInfo("600203", 8, 8, 10000)
	simInfo("600207", 8, 8, 10000)
	simInfo("300001", 3, 3, 10000)
	simInfo("002017", 2, 2, 10000)
	simInfo("603716", 2, 2, 10000)
	//qq.Main()
	InitVal()
	notify.Run()
	golib.Go(jfzt.RunWs)
	golib.Wait()
	log.Info("main exit ")

}
