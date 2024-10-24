package dfcf

import (
	"encoding/json"
	"fmt"
	"io"
	"main/com"
	"main/onebot11"
	"main/redis"
	"net/http"
	"strings"
	"time"
)

type _stockLine struct {
	kai  float64
	shou float64
}

var (
	im com.Notify_xx
)

func init() {
	im = &onebot11.Onebot11Ntf{}
}
func getStockLine(stockId, flag string) []_stockLine {
	/*sz = 0 sh 1*/
	if flag == "sz" {
		flag = "0"
	} else if flag == "sh" {
		flag = "1"
	} else {
		return nil
	}

	url := fmt.Sprintf("https://push2his.eastmoney.com/api/qt/stock/kline/get?"+
		"fields1=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61"+
		"&beg=20240220&end=20500101"+
		"&rtntype=6"+
		"&secid=%s.%s"+
		"&klt=101&fqt=1", flag, stockId)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	m := make(map[string]interface{})
	json.Unmarshal(body, &m)
	data, ok := m["data"].(map[string]interface{})
	if ok == false {
		return nil
	}
	klines, ok := data["klines"].([]interface{})
	if ok == false {
		return nil
	}
	//新股不管
	if len(klines) < 30 {
		return nil
	}
	result := []_stockLine{}
	for _, kline := range klines {
		val, ok := kline.(string)
		if ok == false {
			fmt.Errorf("what kline %v", kline)
			return nil
		}
		valList := strings.Split(val, ",")
		if len(valList) != 11 {
			fmt.Errorf("valList  %v", valList)

			return nil
		}
		x := _stockLine{}
		fmt.Sscanf(valList[1], "%f", &x.kai)
		fmt.Sscanf(valList[2], "%f", &x.shou)
		result = append(result, x)
	}
	return result[len(result)-21:]
}
func GetRa(cur, base float64) float64 {
	return (cur - base) / base * 100
}

// 扫秒热门股
func ScanHotStock() {
	time.AfterFunc(86400*time.Second, func() {
		ScanHotStock()
	})
	data, _ := redis.LoadAll()
	msg := ""

	for k, v := range data {
		info := getStockLine(k, v)
		if info == nil {
			continue
		}
		cntx := 0
		flag := 0
		tmp := 1
		idx := 20
		for idx > 0 {
			if GetRa(info[idx].shou, info[idx-1].shou) >= 9.8 {
				cntx++
				flag |= tmp
			}
			tmp <<= 1
			idx--
		}
		third := 0
		third += flag & 1
		third += (flag & 2) >> 1
		third += (flag & 4) >> 2
		if third > 1 {
			name := redis.Id2Name(k)
			msg = fmt.Sprintf("%s%s 3天%d板 --- 20天%d板\n", msg, name, third, cntx)
		}
	}
	im.SendMsg(msg, nil)

}
