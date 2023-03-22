package ws

import (
	"fmt"
	"main/dingding"
	"main/redis"
	"sync/atomic"
	"time"
)

func (this *WsSet) Ping(stopc chan bool) {
	ping := PingType{ServiceType: "ping"}
	ticker := time.NewTicker(60 * time.Second)
	select {
	case <-ticker.C:
		{
			this.conn.WriteJSON(ping)
		}
	case <-stopc:
		return

	}
}

func (this *WsSet) PostSTATISTICS(id, dq string) {

	dataJson := Data_json{
		SubType:     "SUBON",
		Inst:        id,
		Market:      dq,
		ServiceType: "STATISTICS",
		ReqID:       int(atomic.AddInt32(&reqid, 1)),
	}
	this.conn.WriteJSON(dataJson)

}
func (this *WsSet) PostTick(id, dq string) {
	dataJson := Data_json{
		SubType:     "SUBON",
		Inst:        id,
		Market:      dq,
		ServiceType: "TICK",
		ReqID:       int(atomic.AddInt32(&reqid, 1)),
	}
	this.conn.WriteJSON(dataJson)

}
func (this *WsSet) PostStatic(id, dq string) {

	dataJson := Data_json{
		SubType:     "SUBON",
		Inst:        id,
		Market:      dq,
		ServiceType: "STATIC",
		ReqID:       int(atomic.AddInt32(&reqid, 1)),
	}
	this.conn.WriteJSON(dataJson)

}
func (this *WsSet) PostDyna(id, dq string) {

	dataJson := Data_json{
		SubType:     "SUBON",
		Inst:        id,
		Market:      dq,
		ServiceType: "DYNA",
		ReqID:       int(atomic.AddInt32(&reqid, 1)),
	}
	this.conn.WriteJSON(dataJson)

}

func (this *WsSet) Post(id, dq string) {
	//当前没有post过
	if _, ok := mId2Post[id]; !ok {
		this.PostStatic(id, dq)
		this.PostSTATISTICS(id, dq)
		this.PostTick(id, dq)
		this.PostDyna(id, dq)
		mId2Post[id] = Empty{}
	}
}
func (this *WsSet) PostAndListen(uid, gid, dq string, w int) {
	//当前没有post过
	this.Post(gid, dq)

	AddFollow(gid, uid, w)
}
func PostById(roleid, gid string, w int) {
	id := hx(gid)
	MGR[id].PostAndListen(roleid, gid, redis.GetDQ(gid), w)
}
func GetMgr(id string) *WsSet {
	return MGR[hx(id)]
}
func updtatodd(x bool) {
	msg := ""
	msgzt := ""
	for i := 0; i < WSC; i++ {
		this := MGR[i]
		for k, v := range this.mId2HL {

			if sts, ok := this.mId2BaseData[k]; ok {
				base := sts.PreClosePrice
				bd := GetRa(v.Max, base) - GetRa(v.Min, base)
				if bd > 5 {
					if info, ok := this.mId2ConstInfo[k]; ok {
						msg += fmt.Sprintf("%s   %s 波动较大 %.2f%%\n", info.InstrumentID, info.InstrumentName, bd)
					}
				}
			}

		}
		for k, v := range this.mId2Tick {
			idx := len(v) - 1
			if idx < 0 {
				continue
			}
			x := v[idx]
			base := x.Price
			if sts, ok := this.mId2BaseData[k]; ok {
				base = sts.PreClosePrice
			}
			if info, ok := this.mId2ConstInfo[k]; ok {
				if GetRa(x.Price, base) > 9 && GetRa(x.Price, base) < 50 {
					msgzt += fmt.Sprintf("%s  %s 大涨 %.2f%%\n", info.InstrumentID, info.InstrumentName, GetRa(x.Price, base))

				}
			}
		}
	}
	dingding.SendDingTalkMessage(dingding.DDMsgType{Id: "0", Msg: msgzt}, dingding.KeywordMonitor)
	dingding.SendDingTalkMessage(dingding.DDMsgType{Id: "0", Msg: msg}, dingding.KeywordMonitor)
	if x {
		for _, k := range MGR {
			k.Stop()
		}
	}
	time.AfterFunc(86400*time.Second, func() {
		updtatodd(x)
	})
}
func startListen() {
	for _, x := range MGR {
		x.Reset()
	}
	time.AfterFunc(86400*time.Second, startListen)
}
func daka(msg string) {
	dingding.SendDingTalkMessage(dingding.DDMsgType{Id: "0", Msg: msg}, dingding.KeywordMonitor)

	time.AfterFunc(86400*time.Second, func() {
		daka(msg)
	})
}

// 定时
func DsMsg() {
	{
		timeFormat := "2006-01-02 15:04"
		end, _ := time.ParseInLocation(timeFormat, "2022-04-08 12:31", time.Local)
		diff := time.Now().Sub(end)

		diff %= 86400 * time.Second
		diff = 86400*time.Second - diff

		time.AfterFunc(diff, func() { updtatodd(false) })
	}
	{
		timeFormat := "2006-01-02 15:04"
		end, _ := time.ParseInLocation(timeFormat, "2022-04-08 15:31", time.Local)
		diff := time.Now().Sub(end)

		diff %= 86400 * time.Second
		diff = 86400*time.Second - diff
		time.AfterFunc(diff, func() { updtatodd(true) })

	}
	{
		timeFormat := "2006-01-02 15:04"
		end, _ := time.ParseInLocation(timeFormat, "2022-04-08 09:26", time.Local)
		diff := time.Now().Sub(end)

		diff %= 86400 * time.Second
		diff = 86400*time.Second - diff

		time.AfterFunc(diff, startListen)

	}
	{
		timeFormat := "2006-01-02 15:04"
		end, _ := time.ParseInLocation(timeFormat, "2022-04-08 09:28", time.Local)
		diff := time.Now().Sub(end)

		diff %= 86400 * time.Second
		diff = 86400*time.Second - diff

		time.AfterFunc(diff, func() {
			daka("主人上班别忘记打卡")
		})
	}
	{
		timeFormat := "2006-01-02 15:04"
		end, _ := time.ParseInLocation(timeFormat, "2022-04-08 18:30", time.Local)
		diff := time.Now().Sub(end)

		diff %= 86400 * time.Second
		diff = 86400*time.Second - diff

		time.AfterFunc(diff, func() {
			daka("主人下班别忘记打卡")
		})
	}
}
