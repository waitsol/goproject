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

var DdMsg chan string
var client *http.Client

func init() {
	client = &http.Client{Timeout: 5 * time.Second}
	DdMsg = make(chan string, 100)
}

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
		log.Error(err)
	} else {
		log.Error(resp)
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
		case x := <-DdMsg:
			cache += x
		}
	}
}
