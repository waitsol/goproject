package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// get 方式发起网络请求
func Get(apiURL string, params url.Values) (rs []byte, err error) {
	var Url *url.URL
	Url, err = url.Parse(apiURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return nil, err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	resp, err := http.Get(Url.String())
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

var RedisKey string

func GetByCity(city string) string {
	apiUrl := "http://apis.juhe.cn/simpleWeather/query"

	// 初始化参数
	param := url.Values{}

	// 接口请求参数
	param.Set("city", city)    // 要查询的城市名称/id，城市名称如：温州、上海、北京
	param.Set("key", RedisKey) // 接口请求Key

	// 发送请求
	data, err := Get(apiUrl, param)
	if err != nil {
		// 请求异常，根据自身业务逻辑进行调整修改
		fmt.Errorf("请求异常:\r\n%v", err)
	} else {
		var netReturn map[string]interface{}
		jsonerr := json.Unmarshal(data, &netReturn)
		if jsonerr != nil {
			// 解析JSON异常，根据自身业务逻辑进行调整修改
			fmt.Errorf("请求异常:%v", jsonerr)
		} else {
			errorCode := netReturn["error_code"]
			reason := netReturn["reason"]
			if errorCode.(float64) != 0 {
				fmt.Println(reason)
				return "error"
			}
			data := netReturn["result"]
			// 当前天气信息
			future := data.(map[string]interface{})["future"]

			if errorCode.(float64) == 0 {
				res := future.([]interface{})[0].(map[string]interface{})
				return fmt.Sprintf("天气 %v\n温度 %v", res["weather"], res["temperature"])
			} else {
				// 查询失败，根据自身业务逻辑进行调整修改
				fmt.Printf("请求失败:%v_%v", errorCode.(float64), reason)
			}
		}
	}
	return "error"
}
