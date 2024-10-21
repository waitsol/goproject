package qq

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/waitsol/golib"
	"io"
	"main/notify"
	"net/http"
	"os"
	"time"
)

type Ws struct {
	*websocket.Conn
	s       int
	status  int
	wQueue  chan *Payload
	session string
	url     string
}

func (ws *Ws) write() {
	ticker := time.NewTicker(10 * time.Second)
	for ws.status == 1 {
		select {
		case <-ticker.C:
			return
		case w := <-ws.wQueue:
			ws.WriteJSON(w)
		}
	}

}
func (ws *Ws) heartBeat() {
	payload := &Payload{}
	payload.Op = 0
	for ws.status == 1 {
		payload.D = ws.s
		ws.Conn.WriteJSON(payload)
		time.Sleep(time.Second * 20)
	}
}
func (ws *Ws) stop() {
	ws.status = 0
	ws.Close()
}
func (ws *Ws) recvMsg() {
	for {
		_, b, err := ws.ReadMessage()
		if err != nil {
			log.Error("readmsg error = ", err)
			ws.stop()
			connectWs()
			break
		}
		reply := Payload{}
		json.Unmarshal(b, &reply)
		log.Info("recv ", string(b))
		handleMsg(&reply)
	}
}
func handleMsg(reply *Payload) {
	m, ok := reply.D.(map[string]interface{})

	if ok && m != nil {
		id, _ := m["id"].(string)
		group_id, _ := m["group_id"].(string)
		content, _ := m["content"].(string)
		author, _ := m["author"].(map[string]interface{})
		if author == nil {
			sendMsg("内部错误", "1", "1")
			return
		}
		_, res := notify.GetFollow(author["member_openid"].(string)).HandleMessage(content)
		if res == "" {
			res = "没数据"
		}
		sendMsg(res, id, group_id)
	} else {
		sendMsg("哪里不对劲", "1", "1")
	}
}

func connectWs() *Ws {
	req := newRequest("GET", "/gateway", nil)
	if req == nil {
		log.Error("get ws url req error ")
		return nil
	}

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		log.Error("get ws url error ", err)
		return nil
	}
	data, _ := io.ReadAll(resp.Body)
	m := map[string]string{}
	json.Unmarshal(data, &m)
	wsUrl := m["url"]
	log.Error(string(data))
	dl := websocket.Dialer{}

	header := http.Header{
		"Accept-Language": []string{"zh-CN,zh;q=0.9"},
		"Accept-Encoding": []string{"gzip, deflate, br"},
		"Content-Type":    []string{"application/json"},
		"Authorization":   []string{getToken()},
	}
	conn, _, err := dl.Dial(wsUrl, header)
	if err != nil {
		log.Error("connect ws error ", err)
		return nil
	}
	log.Info("send GateWay ")
	_, b, err := conn.ReadMessage()
	if err != nil {
		log.Error("readmsg error = ", err)
		return nil
	}

	reply := Payload{}
	json.Unmarshal(b, &reply)
	log.Info("recv  GateWay------------", string(b))

	payload := &Payload{}
	payload.Op = 2
	d := make(map[string]interface{}, 30)

	d["token"] = getToken()
	d["intents"] = GROUP_AND_C2C_EVENT
	d["shard"] = []int{0, 1}
	payload.D = d
	//create session
	conn.WriteJSON(payload)

	log.Info("create session ------------------")
	_, b, err = conn.ReadMessage()
	if err != nil {
		log.Error("readmsg error = ", err)
		return nil
	}

	reply = Payload{}
	json.Unmarshal(b, &reply)
	log.Info("recv  session------------", string(b))

	D, ok := reply.D.(map[string]interface{})

	if !ok {
		log.Error("session D error")
		os.Exit(-1)
	}
	session, ok := D["session_id"].(string)
	if !ok || len(session) == 0 {
		log.Error("session  null")
		os.Exit(-1)
	}
	ws := Ws{status: 1, Conn: conn, s: reply.S, session: session}
	ws.wQueue = make(chan *Payload, 20)
	golib.Go(func() {
		ws.heartBeat()
	})
	golib.Go(func() {
		ws.recvMsg()
	})
	return &ws
}
