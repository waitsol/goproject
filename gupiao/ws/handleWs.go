package ws

import (
	"fmt"
	"runtime/debug"
	"time"

	log "github.com/sirupsen/logrus"
)

// 股票异动
func (this *WsSet) handleTick(r dataRes) {
	muban := ""
	v := []int{}
	str := []string{}

	//ramsg := ""
	//bra := false
	if info, ok := this.mId2ConstInfo[r.Inst]; ok {
		muban = info.InstrumentName + "\n"
		//ramsg = fmt.Sprintf("%s  %s  交易量变大\n", r.Inst, info.InstrumentName)
	}
	//标志位 数据是否达到最低发送 波动值
	bflag := false
	//这个是表示是否检查
	bcheck := false
	ch := "↑"
	lastPirce := float64(0)
	//之前成交量数据
	ra, _ok := this.mId2LB[r.Inst]
	if !_ok {
		ra = &VRa{sum: 0}
		this.mId2LB[r.Inst] = ra
	}
	for _, x := range r.QuoteData.TickData {
		bcheck = true
		this.mId2Tick[r.Inst] = append(this.mId2Tick[r.Inst], x)
		//这次一共多少钱
		tprice := x.Price * float64(x.Volume)

		//if this.mId2Time[r.Inst]+20 < x.Time && tprice/ra.GetAvg() >= WarnRatio {
		//	log.WithFields(log.Fields{
		//		"tprice/ra.GetAvg()": tprice / ra.GetAvg(),
		//		"tprice":             tprice,
		//		"ra.GetAvg()":        ra.GetAvg(),
		//		" r.Inst":            r.Inst,
		//	}).Info("tick波动 ")
		//	//	bra = true
		//	this.mId2Time[r.Inst] = x.Time
		//}
		ra.Push(VRaInnner{val: tprice, t: x.Time})

		if x.Price >= lastPirce {
			ch = "↑"
		} else {
			ch = "↓"
		}
		lastPirce = x.Price
		if x.Volume > 200*OneHand {
			base := x.Price
			if sts, ok := this.mId2BaseData[r.Inst]; ok {
				base = sts.PreClosePrice
			} else {
				log.WithField("inst", r.Inst).Error("base empty")
			}
			bflag = true
			if GetRa(x.Price, base) < 22 {
				str = append(str, fmt.Sprintf("%.02f%%   %.02f%s   %d\n", GetRa(x.Price, base), x.Price, ch, x.Volume/OneHand))
				v = append(v, x.Volume)
			} else {
				log.WithFields(log.Fields{
					"x.Price": x.Price, "base": base, "GetRa(x.Price, base)": GetRa(x.Price, base),
				}).Error("tick 计算错误")
			}
		}
	}

	if bflag {
		n := len(v)
		//给关注这个股票的人发消息
		load, ok := SyncId2Listener.Load(r.Inst)
		if ok {
			listens := load.(map[string]int)
			for name, needlen := range listens {
				smsg := muban
				bflag = false
				for i := 0; i < n; i++ {
					if v[i] >= needlen {
						smsg += str[i]
						bflag = true
					}
				}
				if bflag {
					SendMsg(name, smsg)
				}
			}

		}

	}
	if bcheck {
		run := true
		for _, sec := range checkSecs {
			run = this.checkUnActionByTime(r.Inst, sec, run)
		}
		run = true
		for _, cnt := range checkCnts {
			run = this.checkUnActionByCount(r.Inst, cnt, run)
		}

	}
	//
	//if bra {
	//	load, ok := SyncId2Listener.Load(r.Inst)
	//	if ok {
	//		listens := load.(map[string]int)
	//		for name, _ := range listens {
	//			SendMsg(name, ramsg)
	//		}
	//		fmt.Println(bra)
	//		//dingding.DdMsg <- ramsg
	//	}
	//}
}
func SendMsg2Listen(inst, msg string) {

	load, ok := SyncId2Listener.Load(inst)
	if ok {
		listens := load.(map[string]int)
		for name, _ := range listens {
			SendMsg(name, msg)
		}
	}

}
func (this *WsSet) checkUnActionByTime(id string, sec int64, run bool) bool {
	if run == false {
		return false
	}
	now := time.Now().Unix()
	if this.mId2Time[id]+30 < now {
		return false
	}
	//算法 记录最高点和最低点 如果当前点不是最新的最高点和最低点就不上报
	low := float64(60000)
	lidx := -1
	hight := float64(0) //A股不会再有6000点
	hidx := -1
	n := len(this.mId2Tick[id])
	if n <= 0 {
		return false
	}
	endsec := this.mId2Tick[id][n-1].Time
	for i := n - 1; i >= 0; i-- {
		tmp := this.mId2Tick[id][i]
		if endsec-tmp.Time <= sec {
			if tmp.Price >= hight {
				hight = tmp.Price
				hidx = i
			}
			if tmp.Price <= low {
				low = tmp.Price
				lidx = i
			}
		} else {
			break
		}
	}
	if (lidx == n-1 || hidx == n-1) && lidx != hidx {

		muban := ""
		if info, ok := this.mId2ConstInfo[id]; ok {
			muban = info.InstrumentName
		}
		if sts, ok := this.mId2BaseData[id]; ok {
			base := sts.PreClosePrice
			ra1 := GetRa(hight, base)
			ra2 := GetRa(low, base)
			diff := ra1 - ra2
			if diff >= WarnCheck && diff < 20 {
				ch := "↑↑↑"
				if lidx == n-1 {
					ch = "↓↓↓"
				}
				this.mId2Time[id] = now
				SendMsg2Listen(id, fmt.Sprintf("%s 在%d秒异动%s\n%.2f%%  %.2f%%  %.2f%% \n", muban, sec, ch, GetRa(low, base), GetRa(hight, base), diff))
				log.WithFields(log.Fields{
					"lowra": GetRa(low, base), "hightra": GetRa(hight, base), "diff": diff,
					"inst": id, "low": low, "hight": hight, "base": base,
				}).Infof("%s 在%d秒异动 %s", muban, sec, ch)
				return false
			} else if diff < 0 || diff > 20 {
				log.WithFields(
					log.Fields{
						"inst":  id,
						"hight": hight,
						"low":   low,
						"base":  base,
						"diff":  diff,
						"ra1":   ra1,
						"ra2":   ra2,
					}).Error("find checkUnActionByTime catch")
			}
		}
	}
	return true
}
func (this *WsSet) checkUnActionByCount(id string, cnt int, run bool) bool {
	if run == false || this.mId2cnt[id]+60 > time.Now().Unix() {
		return false
	}
	tcnt := cnt
	//算法 记录最高点和最低点 如果当前点不是最新的最高点和最低点就不上报
	low := float64(60000)
	lidx := -1
	hight := float64(0) //A股不会再有6000点
	hidx := -1
	n := len(this.mId2Tick[id])
	if n <= 0 {
		return false
	}
	for i := n - 1; i >= 0 && cnt > 0; i-- {
		cnt--
		tmp := this.mId2Tick[id][i]
		if tmp.Price >= hight {
			hight = tmp.Price
			hidx = i
		}
		if tmp.Price <= low {
			low = tmp.Price
			lidx = i
		}
	}
	if (lidx == n-1 || hidx == n-1) && lidx != hidx {

		muban := ""
		if info, ok := this.mId2ConstInfo[id]; ok {
			muban = info.InstrumentName
		}
		if sts, ok := this.mId2BaseData[id]; ok {
			base := sts.PreClosePrice
			ra1 := GetRa(hight, base)
			ra2 := GetRa(low, base)
			diff := ra1 - ra2
			//fmt.Println("...", hidx, low, diff, base)
			if diff >= WarnRatio && diff < 20 {
				ch := "↑↑↑"
				if lidx == n-1 {
					ch = "↓↓↓"
				}
				this.mId2cnt[id] = time.Now().Unix()
				SendMsg2Listen(id, fmt.Sprintf("%s 在%d次交易中异动 %s\n%.2f%%  %.2f%%  %.2f%% \n", muban, tcnt, ch, GetRa(low, base), GetRa(hight, base), diff))
				log.WithFields(log.Fields{
					"lowra": GetRa(low, base), "hightra": GetRa(hight, base), "diff": diff,
					"inst": id, "low": low, "hight": hight, "base": base,
				}).Infof("%s 在%d次交易中异动 %s", muban, tcnt, ch)
				return false
			} else if diff < 0 || diff > 20 {
				log.WithFields(
					log.Fields{
						"inst":  id,
						"hight": hight,
						"low":   low,
						"base":  base,
						"diff":  diff,
						"ra1":   ra1,
						"ra2":   ra2,
					}).Error("find checkUnActionByTime catch")
			}
		}
	}
	return true
}

// 突破上限才比较
func (this *WsSet) checkUnActionMaxMin(r dataRes) {
	muban := ""
	if info, ok := this.mId2ConstInfo[r.Inst]; ok {
		muban = info.InstrumentName
	}
	if x, ok := this.mId2HL[r.Inst]; ok {
		for _, y := range r.QuoteData.DynaData {
			if y.HighestPrice > x.Max {
				if sts, ok := this.mId2BaseData[r.Inst]; ok {
					ratio1 := GetRa(y.HighestPrice, sts.PreClosePrice)
					ratio2 := GetRa(x.Max, sts.PreClosePrice)
					ratio3 := GetRa(x.Min, sts.PreClosePrice)
					if (ratio1/0.49 != ratio2/0.49) && (ratio1-ratio3 > WarnHL) {
						SendMsg2Listen(r.Inst, fmt.Sprintf("%s 新高↑↑↑\n%.2f%%  %.2f%%\n", muban, ratio1, ratio3))
						log.Infof("%s 新高↑↑↑\n%.2f%%  %.2f%% inst = %s max = %.2f min = %.2f", muban, ratio1, ratio3, r.Inst, x.Max, x.Min)
					}
				}
				x.Max = y.HighestPrice
			}
			if y.LowestPrice < x.Min {
				if sts, ok := this.mId2BaseData[r.Inst]; ok {
					ratio1 := GetRa(y.LowestPrice, sts.PreClosePrice)
					ratio2 := GetRa(x.Min, sts.PreClosePrice)
					ratio3 := GetRa(x.Max, sts.PreClosePrice)

					if (ratio1/0.49 != ratio2/0.49) && (ratio3-ratio1 > WarnHL) {
						SendMsg2Listen(r.Inst, fmt.Sprintf("%s 新低↓↓↓\n%.2f%%  %.2f%%\n", muban, ratio3, ratio1))
						log.Infof("%s 新高↑↑↑\n%.2f%%  %.2f%% inst = %s %.2f %.2f\n", muban, ratio1, ratio3, r.Inst, x.Max, x.Min)

					}
				}
				x.Min = y.LowestPrice
			}
		}
	} else {
		for _, x := range r.QuoteData.DynaData {
			if x.LowestPrice > 0 {
				this.mId2HL[r.Inst] = &HL{Max: x.HighestPrice, Min: x.LowestPrice}
				log.Info("--- ", r.Inst, x.HighestPrice, x.LowestPrice)
			} else {
				log.WithFields(log.Fields{
					"r.Inst": r.Inst, " x.HighestPrice": x.HighestPrice, "x.LowestPrice": x.LowestPrice,
				}).Error("err dyna:")
			}
		}
	}
}

func GetList(ids []string) string {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			debug.PrintStack()
		}
	}()

	msg := "list:"

	for _, f := range ids {
		idx := hx(f)
		if idx < 0 {
			continue
		}
		if sts, ok := MGR[idx].mId2BaseData[f]; ok {
			l := len(MGR[idx].mId2Tick[f])
			if l > 0 {
				tick := MGR[idx].mId2Tick[f][l-1]
				ratio := GetRa(tick.Price, sts.PreClosePrice)
				name := sts.InstrumentID
				if info, ok := MGR[idx].mId2ConstInfo[f]; ok {
					name = info.InstrumentName
				}
				msg += fmt.Sprintf("\n%-10s  %.2f%%", name, ratio)
			}

		}
	}

	return msg
}
func (this *WsSet) handleDyna(r dataRes) {
	for _, x := range r.QuoteData.DynaData {
		this.mId2Dyna[r.Inst] = append(this.mId2Dyna[r.Inst], x)
		this.checkUnActionMaxMin(r)
	}
}
func (this *WsSet) handleSTATISTICS(r dataRes) {
	for _, x := range r.QuoteData.StatisticsData {
		this.mId2BaseData[r.Inst] = x

	}
}
func (this *WsSet) handleStatic(r dataRes) {
	for _, x := range r.QuoteData.StaticData {
		this.mId2ConstInfo[r.Inst] = x
	}
}

func (this *WsSet) handleRes(r dataRes) {

	if this.start && r.ServiceType == "TICK" {
		this.handleTick(r)
	} else if this.start && r.ServiceType == "DYNA" {
		this.handleDyna(r)
	} else if r.ServiceType == "STATISTICS" {
		this.handleSTATISTICS(r)
	} else if r.ServiceType == "STATIC" {
		this.handleStatic(r)
	}
}
