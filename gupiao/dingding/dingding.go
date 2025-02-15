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
	"main/com"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Response struct {
	Errcode int
	Errmsg  string
}

var DDURL string

const KeywordMonitor = ":"

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
type DingDingNtf struct {
}

var id2qq = map[string]string{
	"6DAADAFFBABD5C4CEF7DDDA35F6D1587": "1559556218",
	"4EDC73C71CBA7AD5730606F42BA19204": "529599322",
	"2C7643552C78F85B0A381F23D0213852": "744581755",
}

type POSTMSG struct {
	SenderPlatform string `json:"senderPlatform"`
	ConversationId string `json:"conversationId"`
	AtUsers        []struct {
		DingtalkId string `json:"dingtalkId"`
	} `json:"atUsers"`
	ChatbotCorpId             string `json:"chatbotCorpId"`
	ChatbotUserId             string `json:"chatbotUserId"`
	MsgId                     string `json:"msgId"`
	SenderNick                string `json:"senderNick"`
	IsAdmin                   bool   `json:"isAdmin"`
	SenderStaffId             string `json:"senderStaffId"`
	SessionWebhookExpiredTime int64  `json:"sessionWebhookExpiredTime"`
	CreateAt                  int64  `json:"createAt"`
	SenderCorpId              string `json:"senderCorpId"`
	ConversationType          string `json:"conversationType"`
	SenderId                  string `json:"senderId"`
	ConversationTitle         string `json:"conversationTitle"`
	IsInAtList                bool   `json:"isInAtList"`
	SessionWebhook            string `json:"sessionWebhook"`
	Text                      struct {
		Content string `json:"content"`
	} `json:"text"`
	RobotCode string `json:"robotCode"`
	Msgtype   string `json:"msgtype"`
}

func init() {
	go RecvDDMsg()
}
func (x *DingDingNtf) SendMsg(context string, m map[string]interface{}) {
	id, ok := m["id"].(string)
	if ok { //推送到缓存
		DdMsg <- DDMsgType{Id: id, Msg: context}
	}
}
func init() {
	client = &http.Client{Timeout: 5 * time.Second}
	DdMsg = make(chan DDMsgType, 100)
}

var urlIdx = 0
var urlArr = []string{"https://oapi.dingtalk.com/robot/send?access_token=b428571250dd3e9a82311b26a0001e0a62fcad02690f3d66a6d3c2dc6cae5d31",
	"https://oapi.dingtalk.com/robot/send?access_token=287f75ff42faab4ca8d088e8ab4b94619f7b3409e309bd213b312dc940fd9ac0",
	"https://oapi.dingtalk.com/robot/send?access_token=96c4e0413e790f7650e86514a3d5b0977d84b52d6cc983e4f2c84337b8dcce0e",
	"https://oapi.dingtalk.com/robot/send?access_token=b78271cb4b584719efb411c006b65ebf83ceab99e32e50bf233732e2abf8253d",
	"https://oapi.dingtalk.com/robot/send?access_token=6982ee152bc9709ce56301cb8fc3de7d2b63a31c73e19dbca3b6f360374f85ad"}

// 发送简单文本消息
func SendDingTalkMessage(messageContent []DDMsgType, messagePrefix string) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New("SendDingTalkMessage panic")
		}
	}()
	if com.IsSend() == false {
		//	return nil
	}
	text := map[string]string{}
	at := AtMsg{}
	text["content"] = messagePrefix + ""
	for _, x := range messageContent {

		who := x.Id
		if len(who) == 1 {
			who = "所有人"
		}

		text["content"] = text["content"] + "\n" + x.Msg + "\n" + "@" + who

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
	urlIdx++
	if urlIdx >= len(urlArr) {
		urlIdx = 0

	}
	DDURL = urlArr[urlIdx]
	resp, err := client.Post(DDURL, "application/json", bytes.NewReader(body))
	if err != nil {

		log.Errorf("body = %s,err = %v", string(body), err)
	} else {
		log.Info(resp)
	}
	return nil
}

func RecvDDMsg() { //钉钉缓存
	tick := time.NewTicker(500 * time.Millisecond)
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
