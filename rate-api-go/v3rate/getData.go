package v3rate

import (
	"encoding/json"
	"fmt"
	"rate/api/globe"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var wg sync.WaitGroup

// 按每个国家代码获取汇率 同时更新map中对应的国家汇率字符串
func GetDate(tcur string) {
	defer wg.Done()
	reqs := make(map[string]string)
	reqs["app"] = "finance.rate"
	reqs["scur"] = "CNY"
	reqs["tcur"] = tcur
	reqs["appkey"] = "38501"
	reqs["sign"] = "b4478632262a03a194590a0e555a6914"
	type Result struct {
		Status string `json:"status"`
		Scur   string `json:"scur"`
		Tcur   string `json:"tcur"`
		Ratenm string `json:"ratenm"`
		Rate   string `json:"rate"`
		Update string `json:"update"`
	}
	type Res struct {
		Success string `json:"success"`
		Result  Result `json:"result"`
	}
	res := &Res{}
	rateUrl := "http://api.k780.com"
	client := resty.New()

	_, err := client.R().SetQueryParams(reqs).SetResult(res).ForceContentType("application/json").Get(rateUrl)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	if res.Result.Tcur == "" {
		fmt.Println("未正确获取到汇率:,", tcur)
		return
	}
	mux.Lock()
	defer mux.Unlock()
	//判断globe.CurrencyInfoMap中是否已经存在该货币
	v, ok := globe.CurrencyInfoMap[res.Result.Tcur]
	if ok {
		v.Rate = res.Result.Rate
		globe.CurrencyInfoMap[res.Result.Tcur] = v

	} else {
		fmt.Println("从CF取回数据时发现新货币: ", res.Result.Tcur)

		globe.CurrencyInfoMap[res.Result.Tcur] = globe.CurrencyInfoType{
			Name: res.Result.Tcur + " - " + res.Result.Ratenm,
			Rate: res.Result.Rate,
			// RateNm: res.Result.Ratenm,
			Scur:   res.Result.Scur,
			Tcur:   res.Result.Tcur,
			Update: res.Result.Update,
		}
	}

}

// 当从中国银行获取汇率后，更新v2中的map 以及all_rate_string,然后发送到CF备份
func UpdateRateInfoMapFromApiv2(allRates Rates) {
	//判断globe.All_rate中是否已经存在该货币
	for _, v := range allRates {
		//判断globe.All_rate中是否已经存在该货币
		v := v
		v1, ok := globe.All_rate.Load(v.CurrencyName)
		result := &globe.Result{}
		if ok {
			// fmt.Printf("从cf取回数据中更新rate: %v,原汇率%v,新汇率%v\n", v.CurrencyName, v1.(string), v.Rate)
			err := json.Unmarshal([]byte(v1.(string)), result)
			if err != nil {
				fmt.Printf("err2: %v\n", err)
				return
			}
			result.Rate = v.MiddleRate
			result.Update = strings.Replace(v.PubTime, ".", "-", -1)
			json_bytes, _ := json.Marshal(result)
			globe.All_rate.Store(v.CurrencyName, string(json_bytes))
			fmt.Printf("从v3接口传过来的中国银行数据更新v2的rate: %v,新汇率%v\n", v.CurrencyName, v.MiddleRate)
		} else {
			fmt.Println("从v3接口传过来的中国银行数据中发现新货币: ", v.CurrencyName)
		}

	}
	//调用主包中的函数把map转为字符串

	err := globe.GetRateString()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	//把字符串发送到CF备份

	resps, err := globe.PostRateStringToCF(globe.All_rate_string)
	if err != nil {
		fmt.Printf("err5: %v\n", err)
		return
	}
	fmt.Println(Red, "更新CF中的备份完成", time.Now(), resps, Reset)

}
