/*
-- @Time : 2022/6/8 11:24
-- @Author : raoxiaoya
-- @Desc :
*/
package dingding

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Response struct {
	Errcode int
	Errmsg  string
}

var DDURL string

const KeywordMonitor = "hq"

type DDMsgType struct {
	Id  string
	Msg string
}

var DdMsg chan DDMsgType
var client *http.Client

type AtMsg struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}

func init() {
	client = &http.Client{Timeout: 5 * time.Second}
	DdMsg = make(chan DDMsgType, 100)
}

// 发送简单文本消息
func SendDingTalkMessage(messageContent DDMsgType, messagePrefix string) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New("SendDingTalkMessage panic")
		}
	}()
	who := messageContent.Id
	if len(who) == 1 {
		who = "所有人"
	}
	text := map[string]string{
		"content": messagePrefix + ": " + "@" + who + messageContent.Msg,
	}
	at := AtMsg{}
	at.AtMobiles = []string{messageContent.Id}
	if len(messageContent.Id) == 1 {
		at.IsAtAll = true
	} else {
		at.IsAtAll = false

	}
	postData := map[string]interface{}{
		"msgtype": "text",
		"text":    text,
		"at":      at,
	}

	body, _ := json.Marshal(postData)

	resp, err := client.Post(DDURL, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Error(err)
	} else {
		log.Error(resp)
	}
	return nil
}
func RecvDDMsg() {
	tick := time.NewTicker(4 * time.Second)
	cache := map[string]string{}
	for {
		select {
		case <-tick.C:
			if len(cache) > 0 {
				for k, v := range cache {
					SendDingTalkMessage(DDMsgType{Id: k, Msg: v}, KeywordMonitor)
				}
				cache = map[string]string{}
			}
		case x := <-DdMsg:
			cache[x.Id] += x.Msg
		}
	}
}
