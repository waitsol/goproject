package onebot11

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func init() {
	go loop()
}

var funcQueue = make(chan func(), 20)

type Onebot11Ntf struct {
}

func (x *Onebot11Ntf) SendMsg(context string, m map[string]interface{}) {
	id := m["id"]
	if id != nil && len(id.(string)) > 0 {
		funcQueue <- func() {
			x.sendMsg2group(context, id.(string)) //都发群
		}
	}
	group_id := m["group_id"]
	if group_id != nil && len(group_id.(string)) > 0 {
		x.sendMsg2group(context, group_id.(string))

	}
}
func (*Onebot11Ntf) sendMsg2user(context string, id string) {
	url := "http://127.0.0.1:3000/send_msg"
	body := map[string]interface{}{}
	body["user_id"] = id
	msg := map[string]interface{}{}
	msg["type"] = "text"
	msg["data"] = map[string]string{"text": context}
	body["message"] = msg
	data, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Error(err)
	} else {
		data, _ = io.ReadAll(resp.Body)
		log.Info(string(data))
	}
}
func (*Onebot11Ntf) sendMsg2group(context string, id string) {
	url := "http://127.0.0.1:3000/send_group_msg"
	body := map[string]interface{}{}
	body["group_id"] = 853312133
	msg := map[string]interface{}{}
	msg["type"] = "text"
	msg["data"] = map[string]string{"text": context}
	body["message"] = msg
	data, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Error(err)
	} else {
		data, _ = io.ReadAll(resp.Body)
		log.Info(string(data))
	}
}

func loop() {
	for true {
		f := <-funcQueue
		f()
	}
}
