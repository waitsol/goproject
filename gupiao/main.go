package main

import (
	"main/TcpNet"
	"main/jfzt"
	_ "main/log"
	"main/notify"
	_ "main/onebot11"
	_ "main/redis"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/waitsol/golib"
)

var StartWs chan struct{}

func InitVal() {

	StartWs = make(chan struct{}, 1)

}

func main() {

	go TcpNet.Run()

	time.Sleep(time.Hour)

	//qq.Main()
	InitVal()
	notify.Run()
	golib.Go(jfzt.RunWs)
	golib.Wait()
	log.Info("main exit ")

}
