package main

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/waitsol/golib"
	"main/dingding"
	_ "main/redis"
	"main/wechat"
	"main/ws"
	"time"
)

var StartWs chan struct{}

func InitVal() {

	StartWs = make(chan struct{}, 1)

}
func loginit() {
	logfile := "log/server_"
	linkName := "log/latest_log.log"
	// 日志文件保留的时间

	// 创建日志轮转器
	rotator, err := rotatelogs.New(
		logfile+"%Y%m%d_%H.log",                          // 日志文件名加时间
		rotatelogs.WithLinkName(linkName),                // 始终指向最新的日志文件
		rotatelogs.WithRotationTime(time.Second*60*60*2), //2小时
		rotatelogs.WithMaxAge(time.Second*60*60*24*3),    //3天
	)
	if err != nil {
		log.Fatal("Failed to create rotator: ", err)
	}
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "15:04:05",
	})
	// 设置日志输出到文件
	log.SetOutput(rotator)

	// 打印日志
	log.Info("This is a log message.")
}

func main() {
	loginit()
	InitVal()
	golib.Go(func() {
		wechat.RunWechat(StartWs)
	})

	<-StartWs
	golib.Go(dingding.RecvDDMsg)

	golib.Go(ws.RunWs)
	golib.Wait()
	log.Info("main exit")
}
