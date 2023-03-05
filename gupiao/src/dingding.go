/*
-- @Time : 2022/6/8 11:24
-- @Author : raoxiaoya
-- @Desc :
*/
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	Errcode int
	Errmsg  string
}

var DDURL string

const KeywordMonitor = "hq"

var ddMsg chan string
var client *http.Client

// 发送简单文本消息
func SendDingTalkMessage(messageContent, messagePrefix string) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New("SendDingTalkMessage panic")
		}
	}()
	text := map[string]string{
		"content": messagePrefix + ": " + messageContent,
	}

	postData := map[string]interface{}{
		"msgtype": "text",
		"text":    text,
	}
	body, _ := json.Marshal(postData)

	resp, err := client.Post(DDURL, "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
	return nil
}
func RecvDDMsg() {
	tick := time.NewTicker(4 * time.Second)
	cache := ""
	for {
		select {
		case <-tick.C:
			if len(cache) > 0 {
				SendDingTalkMessage(cache, KeywordMonitor)
				cache = ""
			}
		case x := <-ddMsg:
			cache += x
		}
	}
}
func updtatodd() {
	msg := ""
	msgzt := ""
	for k, v := range mId2HL {

		if sts, ok := mId2BaseData[k]; ok {
			base := sts.PreClosePrice
			bd := getRa(v.Max, base) - getRa(v.Min, base)
			if bd > 5 {
				if info, ok := mId2ConstInfo[k]; ok {
					msg += fmt.Sprintf("%s   %s 波动较大 %.2f%%\n", info.InstrumentID, info.InstrumentName, bd)
				}
			}
		}

	}
	for k, v := range mId2Tick {
		idx := len(v) - 1
		if idx < 0 {
			continue
		}
		x := v[idx]
		base := x.Price
		if sts, ok := mId2BaseData[k]; ok {
			base = sts.PreClosePrice
		}
		if info, ok := mId2ConstInfo[k]; ok {
			if getRa(x.Price, base) > 9 && getRa(x.Price, base) < 50 {
				msgzt += fmt.Sprintf("%s  %s 大涨 %.2f%%\n", info.InstrumentID, info.InstrumentName, getRa(x.Price, base))

			}
		}
	}
	SendDingTalkMessage(msgzt, KeywordMonitor)
	SendDingTalkMessage(msg, KeywordMonitor)
	time.AfterFunc(86400*time.Second, updtatodd)
}
func DsMsg() {
	{
		timeFormat := "2006-01-02 15:04"
		end, _ := time.ParseInLocation(timeFormat, "2022-04-08 12:31", time.Local)
		diff := time.Now().Sub(end)

		diff %= 86400 * time.Second
		diff = 86400*time.Second - diff

		time.AfterFunc(diff, updtatodd)
	}
	{
		timeFormat := "2006-01-02 15:04"
		end, _ := time.ParseInLocation(timeFormat, "2022-04-08 18:31", time.Local)
		diff := time.Now().Sub(end)

		diff %= 86400 * time.Second
		diff = 86400*time.Second - diff

		time.AfterFunc(diff, updtatodd)
	}
	
}
