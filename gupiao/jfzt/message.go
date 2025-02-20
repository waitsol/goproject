package jfzt

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"main/com"
	"main/dfcf"
	"main/dingding"
	"main/redis"
	"sync/atomic"
	"time"
)

var (
	im com.Notify_xx
)

func init() {
	im = &dingding.DingDingNtf{}
}

func (this *WsSet) Ping(stopc chan bool) {
	ping := PingType{ServiceType: "ping"}
	ticker := time.NewTicker(10 * time.Second)
	select {
	case <-ticker.C:
		{
			if this.conn != nil {
				this.conn.WriteJSON(ping)
			}
		}
	case <-stopc:
		log.Info("ping func exit")
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
func (this *WsSet) Listen(uid, gid string, w int) {

	AddFollow(gid, uid, w)
}
func PostById(roleid, gid string, w int) {
	id := hx(gid)
	MGR[id].Listen(roleid, gid, w)
}
func GetMgr(id string) *WsSet {
	return MGR[hx(id)]
}
func calcAvg[T float64](v []T) T {
	ret := T(0)
	for _, x := range v {
		ret += x
	}
	return ret / T(len(v))
}
func checkTurnoverRateByDay(gid string, day int) float64 {
	n := len(mId2TurnoverRate[gid])
	if n <= day {
		return 0
	}
	pre := calcAvg(mId2TurnoverRate[gid][:n-day])
	cur := calcAvg(mId2TurnoverRate[gid][n-day:])
	return GetRa(cur, pre) + 100
}
func GetNameById(this *WsSet, k string) string {
	if info, ok := this.mId2ConstInfo[k]; ok {
		return info.InstrumentName
	}
	return "error"
}
func CheckTurnoverRate() {

	time.AfterFunc(86400*time.Second, func() {
		CheckTurnoverRate()
	})
	if com.IsSend() == false {
		return
	}
	msgzy := "重点:\n"
	msgpt := "普通:\n"
	for i := 0; i < WSC; i++ {
		this := MGR[i]
		for k, v := range this.mId2Dyna {
			n := len(v)
			mId2TurnoverRate[k] = append(mId2TurnoverRate[k], v[n-1].TurnoverRate)
			redis.SaveTurnoverRate(k, v[n-1].TurnoverRate)
			for ra := 500; ra > 350; ra -= 149 {
				for d := 7; d > 0; d-- {
					r := checkTurnoverRateByDay(k, d)
					if int(r) > ra {
						if r >= 500 {
							msgzy += fmt.Sprintf("%v  %v %d天 %.2f%%\n", GetNameById(this, k), k, d, r)
						} else {
							msgpt += fmt.Sprintf("%v  %v %d天 %.2f%%\n", GetNameById(this, k), k, d, r)
						}
						goto exit
					}

				}
			}
		exit:
		}
	}

	if len(msgpt) > 8 {
		im.SendMsg(msgpt, nil)
		time.Sleep(time.Second)

	}
	if len(msgzy) > 8 {
		im.SendMsg(msgzy, nil)
		time.Sleep(time.Second)

	}

}

// 涨停
func updtatodd(x bool) {
	msgzt := ""
	for i := 0; i < WSC; i++ {
		this := MGR[i]
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
			if k[:3] == "300" || k[:3] == "688" {
				continue
			}
			if info, ok := this.mId2ConstInfo[k]; ok {
				if GetRa(x.Price, base) > 9 && GetRa(x.Price, base) < 50 {
					msgzt += fmt.Sprintf("%s  %s 大涨 %.2f%%\n", info.InstrumentID, info.InstrumentName, GetRa(x.Price, base))
				}
			}
		}
	}
	im.SendMsg(msgzt, nil)
	//	dingding.SendDingTalkMessage([]dingding.DDMsgType{{Id: "0", Msg: msgzt}}, dingding.KeywordMonitor)
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
	dingding.SendDingTalkMessage([]dingding.DDMsgType{{Id: "15358698379", Msg: msg}}, dingding.KeywordMonitor)

	time.AfterFunc(86400*time.Second, func() {
		daka(msg)
	})
}

// 定时
func DsMsg() {
	{
		timeFormat := "2006-01-02 15:04"
		end, _ := time.ParseInLocation(timeFormat, "2022-04-08 18:00", time.Local)
		diff := time.Now().Sub(end)

		diff %= 86400 * time.Second
		diff = 86400*time.Second - diff
		time.AfterFunc(diff, func() { updtatodd(true) })

	}
	{
		timeFormat := "2006-01-02 15:04"
		end, _ := time.ParseInLocation(timeFormat, "2022-04-08 09:27", time.Local)
		diff := time.Now().Sub(end)

		diff %= 86400 * time.Second
		diff = 86400*time.Second - diff

		time.AfterFunc(diff, startListen)

	} /*
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
		}*/
	{
		timeFormat := "2006-01-02 15:04"
		end, _ := time.ParseInLocation(timeFormat, "2022-04-08 20:30", time.Local)
		diff := time.Now().Sub(end)

		diff %= 86400 * time.Second
		diff = 86400*time.Second - diff

		time.AfterFunc(diff, func() {
			CheckTurnoverRate()
		})
	}
	{
		timeFormat := "2006-01-02 15:04"
		end, _ := time.ParseInLocation(timeFormat, "2022-04-08 21:30", time.Local)
		diff := time.Now().Sub(end)

		diff %= 86400 * time.Second
		diff = 86400*time.Second - diff

		time.AfterFunc(diff, func() {
			dfcf.ScanHotStock()
		})
	}
}
