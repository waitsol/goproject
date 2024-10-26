package qq

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var AppID = "102452218"
var AccessToken = "8WTKTDyD3IZW1rmfIJpB0WEuZDZVQ5NX"
var url = "https://sandbox.api.sgroup.qq.com"
var AppSecret = "nu19HPXfnv4DMVenw6GQaku4FQbmx8JV"

type Payload struct {
	Op int         `json:"op,omitempty"`
	D  interface{} `json:"d,omitempty"`
	S  int         `json:"s,omitempty"`
	T  string      `json:"t,omitempty"`
}

func getToken() string {

	body := map[string]string{}
	body["appId"] = AppID
	body["clientSecret"] = AppSecret
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return ""
	}
	resp, err := http.Post("https://bots.qq.com/app/getAppAccessToken", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		log.Error("get token error ", err)
		return ""
	}
	data, _ := io.ReadAll(resp.Body)
	m := map[string]string{}
	json.Unmarshal(data, &m)
	log.Error("token ", m["access_token"])

	return fmt.Sprintf("QQBot %s", m["access_token"])
}
func newRequest(method string, path string, reader io.Reader) *http.Request {
	req, _ := http.NewRequest(method, url+path, reader)
	req.Header.Set("Authorization", getToken())
	req.Header.Set("Content-Type", "application/json")
	return req
}

func Main() {

	connectWs()

}

// C2CMessageEventHandler 实现处理 at 消息的回调
//func C2CMessageEventHandler() event.C2CMessageEventHandler {
//	return func(event *dto.WSPayload, data *dto.WSC2CMessageData) error {
//		//TODO use api do sth.
//		return nil
//	}
//}
