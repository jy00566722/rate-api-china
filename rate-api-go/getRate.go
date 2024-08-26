package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"rate/api/globe"
	"rate/api/v3rate"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var wg sync.WaitGroup

// 从CF中取回rate数据,保存入内存中
func InitDateFromCf() error {
	url := "https://rate-back.oeoli.org/get_rate"
	client := resty.New()
	res, err := client.R().Get(url)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	type Result struct {
		Status string `json:"status"`
		Scur   string `json:"scur"`
		Tcur   string `json:"tcur"`
		Ratenm string `json:"ratenm"`
		Rate   string `json:"rate"`
		Update string `json:"update"`
	}
	type RateStringList []string
	var rateStringList RateStringList
	err = json.Unmarshal(res.Body(), &rateStringList)
	if err != nil {
		fmt.Printf("err1: %v\n", err)
		return err
	}
	type RateList []Result
	var rateList RateList
	for _, v := range rateStringList {
		item_result := Result{}
		err = json.Unmarshal([]byte(v), &item_result)
		if err != nil {
			fmt.Printf("err2: %v\n", err)
			return err
		}
		rateList = append(rateList, item_result)
	}
	//保存字符串
	globe.Mux.Lock()
	globe.All_rate_string = string(res.Body())
	globe.Mux.Unlock()
	//保存map
	for _, v := range rateList {
		s, err := json.Marshal(v)
		if err != nil {
			fmt.Printf("err3: %v\n", err)
			return err
		}
		globe.All_rate.Store(v.Tcur, string(s))
	}
	//打印测试保存的值
	// fmt.Println("all_rate:=========")
	// all_rate.Range(func(key, value any) bool {
	// 	fmt.Printf("%v: %v\n", key, value)
	// 	return true
	// })
	// fmt.Println("all_rate_string:=========")
	// fmt.Printf("%v\n", all_rate_string)
	return nil
}

// 定时把每天pop页面的打开数据保存进数据库
func CronSavePopTimesToDb() {
	//增加随机等待时间 500ms到3000ms
	randomTime := rand.Intn(3000) + 500
	time.Sleep(time.Duration(randomTime) * time.Millisecond)

	//连接mongodb
	// var ctx = context.Background()
	// var db, err = qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://rate:rate1234@t.deey.top:57890/?authSource=rate", Database: "rate", Coll: "pop_times"})
	// client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://rate:rate1234@t.deey.top:57890/?authSource=rate"})
	// if err != nil {
	// 	fmt.Printf("连接mongodb出错: %v\n", err)
	// } else {
	// 	fmt.Println("连接mongodb成功-保存pop页面打开数据")
	// }
	// defer client.Close(ctx)
	//总次数
	type PopTimes struct {
		Date  time.Time `bson:"date" json:"date"`
		Times int32     `bson:"times" json:"times"`
	}
	var pop_times = PopTimes{
		Date:  time.Now(),
		Times: PopPageOpenTimes.Load(),
	}
	//按IP的数次
	type PopTimesIp struct {
		Date  time.Time `bson:"date" json:"date"`
		Times int       `bson:"times" json:"times"`
		Ip    string    `bson:"ip" json:"ip"`
	}
	var pop_times_ip []PopTimesIp
	PopPageIpTimes.Range(func(key, value any) bool {
		pop_times_ip = append(pop_times_ip, PopTimesIp{
			Date:  time.Now(),
			Times: value.(int),
			Ip:    key.(string),
		})
		return true
	})
	db := qClient.Database("rate")
	coll1 := db.Collection("pop_times")
	coll2 := db.Collection("pop_times_ip")
	_, err := coll1.InsertOne(context.Background(), pop_times)
	if err != nil {
		fmt.Printf("保存pop页面打开数据出错: %v\n", err)
	} else {
		fmt.Println("保存pop页面打开数据pop_times-成功!")
	}
	_, err = coll2.InsertMany(context.Background(), pop_times_ip)
	if err != nil {
		fmt.Printf("保存pop页面打开数据出错: %v\n", err)
	} else {
		fmt.Println("保存pop页面打开数据pop_times_ip-成功!")
	}
	//清空pop页面打开次数的数据
	PopPageOpenTimes.Store(0)
	PopPageIpTimes.Range(func(key, value any) bool {
		PopPageIpTimes.Delete(key)
		return true
	})
	fmt.Println("清空pop页面打开次数的数据-成功!")
}

// 定时启动获取汇率的具体任务执行
func CtronGetDate(s bool) {
	cur1 := []string{"AFN", "ALL", "AMD", "AOA", "AWG", "AZN", "BAM", "BBD", "BDT", "BGN", "BHD", "BIF", "BMD", "BND", "BOB", "BRL", "BSD", "BTN", "BWP", "BYN", "BYR", "BZD", "CDF", "CRC", "CUC", "CVE", "CZK", "DJF", "DKK"}
	cur2 := []string{"DOP", "DZD", "ERN", "ETB", "FJD", "FKP", "GEL", "GHS", "GIP", "GMD", "GNF", "GTQ", "GYD", "HNL", "HRK", "HTG", "HUF", "IQD", "ISK", "JMD", "JOD", "KES", "KGS", "KHR", "KMF", "KPW", "KWD", "KYD", "KZT", "LAK", "LBP"}
	cur3 := []string{"LRD", "LSL", "LTL", "LVL", "LYD", "MAD", "MDL", "MGA", "MKD", "MMK", "MNT", "MRO", "MRU", "MUR", "MVR", "MWK", "MZN", "NAD", "NGN", "NIO", "NPR", "OMR", "PAB", "PEN", "PGK", "PKR", "PLN", "PYG", "QAR", "RON", "RSD", "RWF"}
	cur4 := []string{"SBD", "SCR", "SDG", "SEK", "SHP", "SLL", "SOS", "SRD", "STD", "SVC", "SYP", "SZL", "TJS", "TMT", "TND", "TOP", "TTD", "TZS", "UGX", "UYU", "UZS", "VEF", "VUV", "WST", "XAF", "XCD", "XDR", "XOF", "XPF", "YER", "ZAR", "ZMW", "ZWL", "ILS"}

	curGood := []string{"CLP", "LKR", "GBP", "COP", "MXN", "MYR", "CAD", "VND", "IRR", "VES", "UAH", "EGP", "CLF", "ARS", "HKD", "NOK", "CUP", "JPY", "CHF", "AUD", "NZD", "USD", "TWD", "CNH", "IDR", "AED", "PHP", "RUB", "MOP", "SAR", "INR", "TRY", "EUR", "THB", "KRW", "SGD", "ANG"}
	// fmt.Printf("curGood: %v\n", curGood)

	d1 := append(cur1, cur2...)
	d2 := append(cur3, cur4...)
	d3 := append(d1, d2...)
	d := append(d3, curGood...) //所有国家代码
	var d_do []string
	if s {
		d_do = d
	} else {
		d_do = curGood
	}
	fmt.Printf("len(d_do): %v\n", len(d_do))
	for _, v := range d_do {
		fmt.Printf("目前在进行的目标: %v\n", v)
		wg.Add(1)
		go getDate(v)
		time.Sleep(time.Millisecond * 150)
	}
	wg.Wait()
	fmt.Println("所有请求完成!!!", time.Now())
	err := globe.GetRateString() //更新内存中的数据
	if err != nil {
		fmt.Printf("err4: %v\n", err)
		return
	}
	fmt.Println("更新内存中的all_rate_string完成", time.Now())

	resps, err := globe.PostRateStringToCF(globe.All_rate_string)
	if err != nil {
		fmt.Printf("err5: %v\n", err)
		return
	}
	fmt.Println("更新CF中的备份完成", time.Now(), resps)

}

// 按每个国家代码获取汇率 同时更新map中对应的国家汇率字符串
func getDate(tcur string) {
	defer wg.Done()
	reqs := make(map[string]string)
	reqs["app"] = "finance.rate"
	reqs["scur"] = "CNY"
	reqs["tcur"] = tcur
	reqs["appkey"] = "38501"
	reqs["sign"] = "b4478632262a03a194590a0e555a6914"

	res := &globe.Res{}
	rateUrl := "http://api.k780.com"
	client := resty.New()

	_, err := client.R().SetQueryParams(reqs).SetResult(res).ForceContentType("application/json").Get(rateUrl)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	json_bytes, _ := json.Marshal(res.Result)

	// fmt.Printf("json_bytes: %v\n", string(json_bytes))
	if res.Result.Tcur == "" {
		fmt.Println("未正确获取到汇率:,", tcur)
		return
	}
	// saveDate(string(res.Result.Tcur), string(json_bytes))
	globe.All_rate.Store(string(res.Result.Tcur), string(json_bytes))

	//传给v3rate,让它更新数据
	v3rate.UpdateCurrencyInfoMapFromApi(res.Result)
}
