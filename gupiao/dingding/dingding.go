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
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
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
func isHoliday(year int, month time.Month, day int) bool {
	// 判断是否是元旦节
	if month == time.January && day == 1 {
		return true
	}

	// 判断是否是春节
	if month == time.February && (day == 11 || day == 12 || day == 13 || day == 14 || day == 15 || day == 16 || day == 17) {
		return true
	}

	// 判断是否是清明节
	if month == time.April && (day == 4 || day == 5 || day == 6) {
		return true
	}

	// 判断是否是劳动节
	if month == time.May && (day == 1 || day == 2 || day == 3) {
		return true
	}

	// 判断是否是端午节
	if month == time.June && (day == 9 || day == 10 || day == 11) {
		return true
	}

	// 判断是否是中秋节
	if month == time.September && (day == 19 || day == 20 || day == 21) {
		return true
	}

	// 判断是否是国庆节
	if month == time.October && (day == 1 || day == 2 || day == 3 || day == 4 || day == 5 || day == 6 || day == 7) {
		return true
	}

	return false
}

func isSend() bool {
	// 获取当前日期
	today := time.Now()

	// 获取今天是周几
	dayOfWeek := int(today.Weekday())

	// 判断今天是否是周六或周日
	if dayOfWeek == 0 || dayOfWeek == 6 {
		return false
	} else {
		// 获取今天的年、月、日
		year, month, day := today.Date()

		// 判断今天是否是法定节假日
		if isHoliday(year, month, day) {
			return false
		} else {
			return true
		}
	}
}

// 发送简单文本消息
func SendDingTalkMessage(messageContent []DDMsgType, messagePrefix string) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New("SendDingTalkMessage panic")
		}
	}()
	if isSend() == false {
		return nil
	}
	text := map[string]string{}
	at := AtMsg{}
	text["content"] = messagePrefix + ":" + "\n"
	for _, x := range messageContent {

		who := x.Id
		if len(who) == 1 {
			who = "所有人"
		}

		text["content"] = text["content"] + "@" + who + x.Msg + "\n"

		at.AtMobiles = append(at.AtMobiles, x.Id)
		if len(x.Id) == 1 {
			at.IsAtAll = true
		} else {
			at.IsAtAll = false
		}

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
				msg := []DDMsgType{}
				for k, v := range cache {
					msg = append(msg, DDMsgType{Id: k, Msg: v})
				}
				SendDingTalkMessage(msg, KeywordMonitor)
				cache = map[string]string{}
			}
		case x := <-DdMsg:
			cache[x.Id] += x.Msg
		}
	}
}
