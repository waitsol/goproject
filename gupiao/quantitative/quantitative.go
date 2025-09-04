package quantitative

import (
	user "main/User"
	"main/ipc"
	"main/proto/go/main/pb3"
	"strconv"

	"github.com/tencent-connect/botgo/log"
)

const (
	S3 = 8.8
	S2 = 5.8
	S1 = 2.8
	B1 = -S1
	B2 = -S2
	B3 = -S3
)

/*
交易限制
1.每只股票不得购买金额不得超过仓位6成
2.量化票最多有2只
3.持仓不得超过总金额的百分之50

系数 东方财富的主力波动决定

交易策略
总购买仓位为6成，卖出2成 购买加一成
跌时情况
1.跌3个点买 1成
2.跌6个点买 2成
3.跌9个点买 3成
急跌情况
急跌3个点以后出现红盘  购买1成 超过3个点 每2个点 加1成
卖出同理
*/

func Buy(code string, ratio float64, money float64, quick_ratio float64) {
	if quick_ratio > 0.5 {
		ratio = quick_ratio
	}
	money += 0.1 //todo后续在算
	stock := user.GetDefaultFollow().FollowsId[code]
	if stock == nil {
		return
	}
	canBuy := stock.BaseHave*15/10 - (stock.All - stock.CanSend)
	canBuy = min(canBuy, stock.Fund/int(money * 100))
	if canBuy <= 0 {
		return
	}
	buyNum := max(stock.BaseHave/10, 1)

	if ratio <= -8.8 {
		buyNum = min(buyNum*3, canBuy)
	} else if ratio <= -6 {
		buyNum = min(buyNum*2, canBuy)
	}
	buyNum *= 100
	Num := strconv.Itoa(buyNum)
	Money := strconv.FormatFloat(money, 'f', 2, 64)
	buy := true
	bs := &pb3.CSBuyOrSellReq{Code: &code, Num: &Num, Money: &Money, IsBuy: &buy}
	req := &pb3.PacketReq{Packet: &pb3.PacketReq_Bs{Bs: bs}}
	log.Info("......zzBuy.....", req.String())
	ipc.ReqChan <- *req
	//更新当前股票资金 拥有数量
	stock.Fund -= buyNum * int(money)
	stock.All += buyNum / 100
}
func Sell(code string, ratio float64, money float64, quick_ratio float64) {
	if quick_ratio > 0.5 {
		ratio = quick_ratio
	}
	money -= 0.1
	stock := user.GetDefaultFollow().FollowsId[code]
	if stock == nil {
		return
	}
	canSell := stock.CanSend
	if canSell <= 0 {
		return
	}

	sellNum := max(stock.BaseHave/10, 1)

	if ratio >= 8.8 {
		sellNum = min(sellNum*3, canSell)
	} else if ratio >= 5.8 {
		sellNum = min(sellNum*2, canSell)
	}

	sellNum *= 100
	Num := strconv.Itoa(sellNum)
	Money := strconv.FormatFloat(money, 'f', 2, 64)
	buy := true
	bs := &pb3.CSBuyOrSellReq{Code: &code, Num: &Num, Money: &Money, IsBuy: &buy}
	req := &pb3.PacketReq{Packet: &pb3.PacketReq_Bs{Bs: bs}}
	log.Info("......zzsell.....", req.String())
	ipc.ReqChan <- *req
	stock.Fund -= sellNum * int(money)
	//更新能卖的数量
	stock.CanSend -= sellNum / 100
	stock.All -= sellNum / 100
	stock.Fund += sellNum * int(money)

}
