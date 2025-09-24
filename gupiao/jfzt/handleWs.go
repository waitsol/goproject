package jfzt

import (
	"fmt"
	user "main/User"
	"main/quantitative"
	"runtime/debug"
	"time"

	log "github.com/sirupsen/logrus"
)

func (this *WsSet) Stock(inst string) *Stock {
	stock := this.stock[inst]
	if stock == nil {
		stock = &Stock{
			hl:        HL{},
			baseData:  StatisticType{},
			constInfo: StaticType{},
			tick:      make([]TickType, 0, 50), // 建议预分配合理容量
			dyna:      make([]DynaType, 0, 50),
			lb:        &VRa{sum: 0}, // 如果 VRa 也需要初始化
			up:        KvInfo{0, 0, 0},
			down:      KvInfo{0, 0, 0},
			buy:       TsInfo{},
			sell:      TsInfo{},
			msFlag:    make([]bool, 10),
		}
		this.stock[inst] = stock
	}
	stock.msFlag[0] = true
	return stock
}

// 股票异动
func (this *WsSet) handleTick(r dataRes) {
	muban := ""
	v := []int{}
	str := []string{}
	stock := this.Stock(r.Inst)
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
	if len(this.mId2Tick[r.Inst]) > 0 {
		lastPirce = this.mId2Tick[r.Inst][len(this.mId2Tick[r.Inst])-1].Price
	}
	if r.Inst != "601606" {
		return
	}
	//之前成交量数据
	ra, _ok := this.mId2LB[r.Inst]
	if !_ok {
		ra = &VRa{sum: 0}
		this.mId2LB[r.Inst] = ra
	}
	for _, x := range r.QuoteData.TickData {
		bcheck = true
		this.mId2Tick[r.Inst] = append(this.mId2Tick[r.Inst], x)

		base := x.Price
		if sts, ok := this.mId2BaseData[r.Inst]; ok {
			base = sts.PreClosePrice
		} else {
			log.WithField("inst", r.Inst).Error("base empty")
		}
		bigSell := time.Now().Unix() - stock.bigSell

		bigBuy := time.Now().Unix() - stock.bigBuy

		ratio := GetRa(x.Price, base)
		msIdx := 0
		if ratio >= quantitative.S3 {
			msIdx = 1
		} else if ratio >= quantitative.S2 {
			msIdx = 2

		} else if ratio >= quantitative.S1 {
			msIdx = 3

		} else if ratio <= quantitative.B3 {
			msIdx = 4

		} else if ratio <= quantitative.B2 {
			msIdx = 5

		} else if ratio <= quantitative.B1 {
			msIdx = 6
		}
		if !stock.msFlag[msIdx] {
			stock.msFlag[msIdx] = true
			if ratio > 0 && bigBuy > 10 {
				log.Infof("%v zz Sell ,ratio=%.2f ", r.Inst, ratio)
				quantitative.Sell(r.Inst, ratio, x.Price, 0)
			} else if bigSell > 10 {
				log.Infof("%v zz Buy ,ratio=%.2f ", r.Inst, ratio)

				quantitative.Buy(r.Inst, ratio, x.Price, 0)
			}
		}

		if x.Price >= lastPirce {
			ch = "↑"
			stock.up.cnt++
			stock.up.val += int64(x.Price) * int64(x.Volume)
			if stock.up.cnt == 1 {
				stock.up.begin = x.Price
			}
			if stock.up.cnt == 2 {
				//抄底买
				xRa := GetRa(stock.down.begin, base)
				if xRa-ratio > 3 {
					log.Infof("%v zz buy stock.down.cnt = %d,ratio=%.2f xRa = %.2f", r.Inst, stock.down.cnt, ratio, xRa)
					quantitative.Buy(r.Inst, ratio, x.Price, xRa)
				}
				stock.down.Reset()
			}
		} else {
			ch = "↓"
			stock.down.cnt++
			stock.down.val += int64(x.Price) * int64(x.Volume)
			if stock.down.cnt == 1 {
				stock.down.begin = x.Price
			}
			if stock.down.cnt == 2 {
				//拉高回调卖
				xRa := GetRa(stock.up.begin, base)
				if ratio-xRa > 3 {
					log.Infof("%v zz sell stock.up.cnt = %d,ratio=%.2f xRa = %.2f", r.Inst, stock.up.cnt, ratio, xRa)
					quantitative.Sell(r.Inst, ratio, x.Price, xRa)
				}
				stock.up.Reset()
			}
		}
		lastPirce = x.Price
		xishu := int(max(base/10, 1)) //10块钱是一万手算大单

		if x.Volume > 10000*OneHand/xishu {
			if ch == "↓" {
				stock.bigSell = time.Now().Unix()
			} else {
				stock.bigBuy = time.Now().Unix()
			}
			bflag = true
			if ratio < 22 {
				str = append(str, fmt.Sprintf("%.02f%%   %.02f%s   %d\n", GetRa(x.Price, base), x.Price, ch, x.Volume/OneHand))
				//log.Info(str)
				v = append(v, x.Volume)
			} else {
				log.WithFields(log.Fields{
					"x.Price": x.Price, "base": base, "GetRa(x.Price, base)": GetRa(x.Price, base),
				}).Error("tick 计算错误")
			}
		} else {
		}
	}

	if bflag {
		n := len(v)
		//给关注这个股票的人发消息
		load, ok := user.SyncId2Listener.Load(r.Inst)
		if ok {
			listens := load.(map[string]*user.FollowInfo)
			bF := true
			for name, info := range listens {
				smsg := muban
				bflag = false
				for i := 0; i < n; i++ {
					if v[i] >= info.Num*OneHand {
						smsg += str[i]
						bflag = true
					}
				}
				if bflag {
					SendMsg(name, smsg, bF)
					bF = false
				}
			}

		}
	}
	if bcheck {
		run := true
		for _, cnt := range checkCnts {
			run = this.checkUnActionByCount(r.Inst, cnt, run)
		}
		for _, sec := range checkSecs {
			run = this.checkUnActionByTime(r.Inst, sec, run)
		}
	}

	if sts, ok := this.mId2BaseData[r.Inst]; ok {
		base := sts.PreClosePrice
		curra := GetRa(lastPirce, base)
		load, ok := user.SyncId2Listener.Load(r.Inst)
		if ok {
			listens := load.(map[string]*user.FollowInfo)

			for id := range listens {
				x := listens[id]
				if x != nil {
					if curra >= x.MaxRa {
						SendMsg(id, fmt.Sprintf("%s %f 超过提醒值 %f ", muban, x.MaxRa, curra), false)
						x.MaxRa = 200 //就提醒一次
					}
					if curra <= x.MinRa {
						SendMsg(id, fmt.Sprintf("%s %f 超过提醒值 %f ", muban, x.MinRa, curra), false)
						x.MinRa = -200
					}
				}
			}
		}
	}
	//
	//if bra {
	//	load, ok := user.SyncId2Listener.Load(r.Inst)
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
	load, ok := user.SyncId2Listener.Load(inst)
	if ok {
		listens := load.(map[string]*user.FollowInfo)
		flag := true
		for name := range listens {
			SendMsg(name, msg, flag)
			flag = false
		}
	}

}
func (this *WsSet) checkUnActionByTime(id string, sec int64, run bool) bool {
	if run == false {
		return false
	}
	now := time.Now().Unix()
	if this.mId2Time[id]+30 > now {
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
						log.Infof("%s 新低↓↓↓\n%.2f%%  %.2f%% inst = %s %.2f %.2f\n", muban, ratio1, ratio3, r.Inst, x.Max, x.Min)

					}
				}
				x.Min = y.LowestPrice
			}
		}
	} else {
		for _, x := range r.QuoteData.DynaData {
			if x.LowestPrice > 0 {
				this.mId2HL[r.Inst] = &HL{Max: x.HighestPrice, Min: x.LowestPrice}
				//log.Info("--- ", r.Inst, x.HighestPrice, x.LowestPrice)
			} else {
				//log.WithFields(log.Fields{
				//	"r.Inst": r.Inst, " x.HighestPrice": x.HighestPrice, "x.LowestPrice": x.LowestPrice,
				//}).Warn("err dyna:")
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

	if r.ServiceType == "TICK" {
		if MGR[0].start {
			this.handleTick(r)
		}
	} else if r.ServiceType == "DYNA" {
		this.handleDyna(r)
	} else if r.ServiceType == "STATISTICS" {
		this.handleSTATISTICS(r)
	} else if r.ServiceType == "STATIC" {
		this.handleStatic(r)
	}

}
