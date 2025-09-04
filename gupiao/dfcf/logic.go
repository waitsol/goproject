package dfcf

import (
	"fmt"
	"main/onebot11"
	"main/redis"
	"sort"
	"strconv"
	"time"

	"github.com/waitsol/golib"
)

var kaipan = map[string]float64{}
var onebot = onebot11.Onebot11Ntf{}

func init() {
	// golib.Go(func() {
	// 	KaiPanIng()
	// })
	// golib.Go(
	// 	func() {
	// 		tingPan()
	// 	},
	// )
}
func KaiPanIng() {
	time.Sleep(time.Hour)
	golib.Go(func() {
		timeFormat := "2006-01-02 15:04"
		end, _ := time.ParseInLocation(timeFormat, "2022-04-08 18:00", time.Local)
		diff := time.Now().Sub(end)

		diff %= 86400 * time.Second
		diff = 86400*time.Second - diff
		time.Sleep(diff)
		for {
			kaipan = map[string]float64{}
			time.Sleep(time.Hour * 24)
		}
	})

	for {
		time.Sleep(time.Second * 30)
		result := getgailiangu()
		begin := 9*60 + 30
		end := 15*60 + 30
		now := time.Now().Hour()*60 + time.Now().Second()
		if now < begin || now > end {
			continue
		}

		ntf := ""
		for x, info := range result {
			if info.zf-kaipan[x] > 0.5 {
				ntf = fmt.Sprintf("%s\n%s概念快速拉伸 %.2f\n", ntf, x, info.zf-kaipan[x])
			}
			kaipan[x] = info.zf
		}
		onebot.SendMsg(ntf, nil)
	}

}

func tingPan() {
	golib.Go(func() {
		timeFormat := "2006-01-02 15:04"
		end, _ := time.ParseInLocation(timeFormat, "2022-04-08 12:00", time.Local)
		diff := time.Now().Sub(end)

		diff %= 86400 * time.Second
		diff = 86400*time.Second - diff
		time.Sleep(diff)
		for {
			kaipan = map[string]float64{}
			time.Sleep(time.Hour * 24)
		}
	})

	timeFormat := "2006-01-02 15:04"
	end, _ := time.ParseInLocation(timeFormat, "2022-04-08 22:00", time.Local)
	diff := time.Now().Sub(end)

	diff %= 86400 * time.Second
	diff = 86400*time.Second - diff
	time.Sleep(diff)
	for {
		result := getgailiangu()
		//什么板块大涨
		ntf := "前30涨幅概念"
		vec := []concept{}
		for _, x := range result {
			vec = append(vec, x)
		}
		sort.SliceStable(vec, func(i, j int) bool {
			return vec[i].zf > vec[j].zf
		})
		for i := 0; i < 30; i++ {
			ntf = fmt.Sprintf("%s\n%s概念 : %.2f", ntf, vec[i].name, vec[i].zf)
		}
		onebot.SendMsg(ntf, nil)
		//什么板块涨停最多
		sort.SliceStable(vec, func(i, j int) bool {
			return vec[i].zt > vec[j].dz
		})
		ntf = ""
		for i := 0; i < 20; i++ {
			ntf = fmt.Sprintf("%s\n%s概念涨停%d个...相关信息%s\n", ntf, vec[i].name, vec[i].zt, vec[i].ntfInfo)
		}
		onebot.SendMsg(ntf, nil)
		time.Sleep(time.Hour * 24)
	}
}

// 获取股票资金情况 近4分钟的情况,当前值,最大值,最小值
func GetStockFundInfo(stockId string) (int, int, int, int) {
	vec := PullFundInfo(stockId, redis.GetDQ(stockId))
	if len(vec) == 0 {
		return -1000, 0, 0, 0
	}
	cur := 0
	vmax := -0x3f3f3f3f3f
	vmin := 0x3f3f3f3f3f
	qushi := 0

	for _, x := range vec {
		vf, _ := strconv.ParseFloat(x, 64)
		vx := int(vf)
		if vx > vmax {
			vmax = vx
		}
		if vx < vmin {
			vmin = vx
		}
		cur = vx
	}

	qmax := -0x3f3f3f3f3f
	qmin := 0x3f3f3f3f3f

	for i := max(0, len(vec)-15); i < len(vec); i++ {
		vf, _ := strconv.ParseFloat(vec[i], 64)
		vx := int(vf)
		if vx > qmax {
			qmax = vx
		}
		if vx < qmin {
			qmin = vx
		}
	}
	if qmax-cur > cur-qmin {
		qushi = cur - qmax
	} else {
		qushi = cur - qmin
	}

	return qushi / 10000, cur / 10000, vmax / 10000, vmin / 10000
}
