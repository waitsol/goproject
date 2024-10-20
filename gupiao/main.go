package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/waitsol/golib"
	_ "main/log"
	"main/notify"
	_ "main/onebot11"
	"main/qq"
	_ "main/redis"
	"main/ws"
)

var StartWs chan struct{}

func InitVal() {

	StartWs = make(chan struct{}, 1)

}

func main() {

	qq.Main()
	InitVal()
	notify.Run()

	golib.Go(ws.RunWs)
	golib.Wait()
	log.Info("main exit ")

}
