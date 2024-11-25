package dfcf

import (
	"fmt"
	"github.com/waitsol/golib"
	"main/onebot11"
	"sort"
	"time"
)

var kaipan = map[string]float64{}
var onebot = onebot11.Onebot11Ntf{}

func init() {
	golib.Go(func() {
		KaiPanIng()
	})
	golib.Go(
		func() {
			tingPan()
		},
	)
}
func KaiPanIng() {
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
