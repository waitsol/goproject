package qq

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func sendMsg(msg, id, group_id string) {
	m := map[string]interface{}{}
	m["content"] = msg
	m["msg_type"] = 0
	if len(id) > 0 {
		m["msg_id"] = id
	}
	data, _ := json.Marshal(m)

	req := newRequest("POST", "/v2/groups/"+group_id+"/messages", bytes.NewReader(data))
	log.Infof("send %v", m)
	//req.Body.Read(data)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("send msg data = %v error %v", msg, err)
		return
	}
	data, _ = io.ReadAll(resp.Body)
	log.Info("recv resp ", string(data))
	json.Unmarshal(data, &m)

}
