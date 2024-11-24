package dfcf

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
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
func getConceptAll(f12 string, pn int) map[string]float64 {

	result := map[string]float64{}
	total := 21
	for pn*20 < total {
		url := fmt.Sprintf("https://89.push2.eastmoney.com/api/qt/clist/get?&pn=%d&pz=20&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&dect=1&wbp2u=|0|0|0|web&fid=f3&fs=b:%s+f:!50&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152,f45&_=1732458489443", pn, f12)
		pn++
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
		tt, ok := data["total"].(float64)
		if ok == false {
			return nil
		}
		total = int(tt)
		diffs, ok := data["diff"].([]interface{})
		if ok == false {
			return nil
		}
		for _, diff := range diffs {
			tmp, ok := diff.(map[string]interface{})
			if ok == false {
				return nil
			}
			f14, ok := tmp["f14"].(string) //股票名字 f12 是id string
			if ok == false {
				return nil
			}
			f3, ok := tmp["f3"].(float64) //涨跌幅度
			if ok == false {
				return nil
			}
			log.Debug("%s : %f\n", f14, f3)
			result[f14] = f3
		}
	}

	return result
}

type concept struct {
	zf      float64
	childs  map[string]float64
	name    string
	zt      int //涨停
	dz      int //大涨
	ntfInfo string
}

func getgailiangu() map[string]concept {
	ret := map[string]concept{}
	pn := 1
	total := 21
	for pn*20 < total {
		url := fmt.Sprintf("https://53.push2.eastmoney.com/api/qt/clist/get?cb=jQuery11240837685905576268_1732454200402&pn=%d&pz=20&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&dect=1&wbp2u=|0|0|0|web&fid=f3&fs=m:90+t:3+f:!50&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f26,f22,f33,f11,f62,f128,f136,f115,f152,f124,f107,f104,f105,f140,f141,f207,f208,f209,f222&_=1732453201149", pn)
		pn++
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		body, _ := io.ReadAll(resp.Body)
		m := make(map[string]interface{})
		if len(body) < 300 {
			return nil
		}
		bidx := strings.Index(string(body), "{")
		eidx := strings.LastIndex(string(body), "}")
		if bidx == -1 || eidx == -1 {
			return nil
		}
		bb := body[bidx : eidx+1]
		//fmt.Println(string(bb))
		json.Unmarshal(bb, &m)
		data, ok := m["data"].(map[string]interface{})
		if ok == false {
			return nil
		}
		tt, ok := data["total"].(float64)
		if ok == false {
			return nil
		}
		total = int(tt)
		//处理数据
		diffs, ok := data["diff"].([]interface{})
		if ok == false {
			return nil
		}
		for _, diff := range diffs {
			tmp, ok := diff.(map[string]interface{})
			if ok == false {
				return nil
			}
			f12, ok := tmp["f12"].(string) //概念的编号
			if ok == false {
				return nil
			}
			f3, ok := tmp["f3"].(float64) //粘贴幅度
			if ok == false {
				return nil
			}
			f14, ok := tmp["f14"].(string)
			if ok == false {
				return nil
			}
			log.Debug("概念:", f14)
			result := getConceptAll(f12, 1)
			zt := 0
			dz := 0
			ntfInfo := ""
			for name, x := range result {
				if x > 9.85 {
					zt++
					ntfInfo = fmt.Sprintf("%s\n%s:%.2f", ntfInfo, name, x)
				}
				if x > 8 {
					dz++
				}
			}
			ret[f14] = concept{f3, result, f14, zt, dz, ntfInfo}

		}
	}

	return ret
}
