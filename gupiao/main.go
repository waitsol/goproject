package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/waitsol/golib"
	"main/jfzt"
	_ "main/log"
	"main/notify"
	_ "main/onebot11"
	"main/qq"
	_ "main/redis"
)

var StartWs chan struct{}

func InitVal() {

	StartWs = make(chan struct{}, 1)

}

func main() {

	qq.Main()
	InitVal()
	notify.Run()

	golib.Go(jfzt.RunWs)
	golib.Wait()
	log.Info("main exit ")

}
