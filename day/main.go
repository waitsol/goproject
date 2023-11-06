package main

import (
	"day/dingding"
	_ "day/redis"
	"day/weather"
	"time"
)

func TimerFunc(t string, f func()) {
	timeFormat := "2006-01-02 15:04"
	end, _ := time.ParseInLocation(timeFormat, t, time.Local)
	diff := time.Now().Sub(end)

	diff %= 86400 * time.Second
	diff = 86400*time.Second - diff
	time.AfterFunc(diff, func() {
		f()
		TimerFunc(t, f)
	})
}

func main() {
	TimerFunc("2006-01-02 19:00", func() {
		dingding.SendDingTalkMessage([]dingding.DDMsgType{{Id: `15358698379`, Msg: weather.GetByCity("荆州")}}, dingding.KeywordMonitor)
	})
	time.Sleep(time.Hour * 999999)
}
