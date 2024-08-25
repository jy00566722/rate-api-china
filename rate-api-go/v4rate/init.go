package v4rate

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rate/api/globe"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type CurrencyInfo globe.CurrencyInfoType

// 从cloudflare的api获取备份的汇率 Key:all-rate-v4
func InitDataFromCFV4() {
	fmt.Println("从cf备份中获取数据v4,更新全局变量!!!")
	url := "https://rate-back.oeoli.org/get_rate_v4"
	client := resty.New()
	res, err := client.R().Get(url)
	if err != nil {
		fmt.Printf("请求cf中的v4备份数据出错,err: %v\n", err)
		panic(err)
	}
	var rateMap map[string]CurrencyInfo
	err = json.Unmarshal(res.Body(), &rateMap)
	if err != nil {
		fmt.Printf("解析cf中回来的v4数据出错,err: %v\n", err)
		panic(err)
	}
	for k, v := range rateMap {
		newUpdate := "c" + v.Update[1:]
		v.Update = newUpdate //把cf取回来的值的update字段更新一下，方便识别
		globe.V4CruuencyMap.Store(k, v)
	}
}

// 从 now-api中请求很少见的货币汇率-这些汇率在google sheets中没有
func GetRateFromNowApi() {
	fmt.Println("从now-api中获取数据")
	cu_list := []string{"CUC", "BYR", "ERN", "FKP", "GIP", "HRK", "KPW", "LTL", "LVL", "MNT", "MRO", "SHP", "STD", "SYP", "VEF", "VUV", "WST", "XDR", "ZWL"}
	// cu_list := []string{"ERN", "FKP", "GIP", "HRK", "KPW", "MNT", "SHP", "SYP", "VUV", "WST", "ZWL"}
	for _, v := range cu_list {
		time.Sleep(150 * time.Millisecond)
		res, err := getDate(v)
		if err != nil {
			fmt.Println("未正确获取到汇率:,", v)
		} else {
			// globe.V4CruuencyMap.Store(res.Result.Tcur, res)
			v1, ok := globe.V4CruuencyMap.Load(res.Result.Tcur)
			if ok {
				v2, ok := v1.(CurrencyInfo)
				if ok {
					// t := v2.Rate
					v2.Rate = res.Result.Rate
					v2.Update = "n" + getFormattedTime()
					globe.V4CruuencyMap.Store(res.Result.Tcur, v2)
					// fmt.Println("从now-api更新汇率:", res.Result.Tcur, "原:", t, "新:", res.Result.Rate)
				} else {
					fmt.Println("类型断言失败:", v)
				}
			} else {
				fmt.Println("从now-api获取数据时发现新货币,未进一步处理:", v)
			}
		}
	}
	fmt.Println("从now-api中获取数据-完毕")
}

// 从 now-api中请求常见的货币汇率
func GetRateFromNowApiCommon() {
	fmt.Println("从now-api中获取数据-常见货币")
	cu_list := []string{
		"USD",
		"EUR",
		"GBP",
		"JPY",
		"AUD",
		"KRW",
		"CAD",
		"NZD",
		"MOP",
		"HKD",
		"TWD",
		"AED",
		"SAR",
		"BRL",
		"CHF",
		"DKK",
		"SEK",
		"NOK",
		"IDR",
		"MYR",
		"PHP",
		"SGD",
		"THB",
		"INR",
		"RUB",
		"TRY",
		"ZAR",
	}
	// cu_list := []string{"ERN", "FKP", "GIP", "HRK", "KPW", "MNT", "SHP", "SYP", "VUV", "WST", "ZWL"}
	for _, v := range cu_list {
		time.Sleep(150 * time.Millisecond)
		res, err := getDate(v)
		if err != nil {
			fmt.Println("未正确获取到汇率:,", v)
		} else {
			// globe.V4CruuencyMap.Store(res.Result.Tcur, res)
			v1, ok := globe.V4CruuencyMap.Load(res.Result.Tcur)
			if ok {
				v2, ok := v1.(CurrencyInfo)
				if ok {
					// t := v2.Rate
					v2.Rate = res.Result.Rate
					v2.Update = "n" + getFormattedTime()
					globe.V4CruuencyMap.Store(res.Result.Tcur, v2)
					// fmt.Println("从now-api更新汇率:", res.Result.Tcur, "原:", t, "新:", res.Result.Rate)
				} else {
					fmt.Println("类型断言失败:", v)
				}
			} else {
				fmt.Println("从now-api获取数据时发现新货币,未进一步处理:", v)
			}
		}
	}
	fmt.Println("从now-api中获取数据-常见货币-完毕")
}

// 从google sheets中获取汇率
func GetRateFromGoogleSheets() {
	fmt.Println("从google sheets中获取数据-开始")
	// 发送 HTTP GET 请求
	resp, err := http.Get("https://docs.google.com/spreadsheets/d/e/2PACX-1vQS4NljENMDBH7xGLQ7OpmogFqm7mnwMg_W6MctPGxfaC0GJbitgFejoRUCAlmo9fh1k75DfXhOeMRN/pub?gid=0&single=true&output=csv")
	if err != nil {
		fmt.Println("Error fetching the CSV file:", err)
		return
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: HTTP status code", resp.StatusCode)
		return
	}

	// 创建 CSV 读取器
	reader := csv.NewReader(resp.Body)

	// 读取所有行
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV data:", err)
		return
	}
	for _, record := range records {
		rate := record[0]
		code := record[1]
		//这里修复一个bug,rate不能是无效的汇率数字字符串,否则会导致汇率不正常,客户端出错
		_, err := strconv.ParseFloat(rate, 64)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		v, ok := globe.V4CruuencyMap.Load(code)
		if ok {
			existing, ok := v.(CurrencyInfo)
			if ok {
				// fmt.Println("从google sheets更新汇率:", code, "原:", existing.Rate, "新:", rate)
				existing.Rate = rate
				existing.Update = "g" + getFormattedTime()
				globe.V4CruuencyMap.Store(code, existing)

			} else {
				fmt.Println("从google sheets获取数据时类型断言失败:", code)
			}
		} else {
			fmt.Println("从google sheets获取数据时发现新货币,未进一步处理:", code)
		}

	}
	fmt.Println("从google sheets中获取数据-结束")

}

// 从google sheets中获取相对美元的汇率，然后转换为相对人民币的汇率-这是一个应急的功能，不能常用
func GetRateFromGoogleSheetsUSD() {
	// 发送 HTTP GET 请求
	resp, err := http.Get("https://docs.google.com/spreadsheets/d/e/2PACX-1vQlmuo1VKRfqHEmzRslTHf1xxAukaMo8TK3kNM7FUGqAhfjH16RoqUb0gJOf71oOzEfaJhQVx6SsyON/pub?gid=0&single=true&output=csv")
	if err != nil {
		fmt.Println("Error fetching the CSV file:", err)
		return
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: HTTP status code", resp.StatusCode)
		return
	}

	// 创建 CSV 读取器
	reader := csv.NewReader(resp.Body)

	// 读取所有行
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV data:", err)
		return
	}
	for _, record := range records {
		rate := record[0]
		code := record[1]
		v, ok := globe.V4CruuencyMap.Load(code)
		if ok {
			existing, ok := v.(CurrencyInfo)
			if ok {
				// fmt.Println("从google sheets更新汇率:", code, "原:", existing.Rate, "新:", rate)
				rate_num, err := strconv.ParseFloat(rate, 64)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				existing.Rate = fmt.Sprint(rate_num / 7.16822)
				existing.Update = "g" + getFormattedTime()
				globe.V4CruuencyMap.Store(code, existing)

			} else {
				fmt.Println("从google sheets获取数据时类型断言失败:", code)
			}
		} else {
			fmt.Println("从google sheets获取数据时发现新货币,未进一步处理:", code)
		}

	}

}

// 从er-api获取汇率 这个api每天早上8点更新(东八区)
func GetRateFromErApi() {
	data, err := getExchangData()
	if err != nil {
		fmt.Println("从er-api获取汇率时出错:", err)
		return
	}
	cu_list := []string{
		"USD",
		"EUR",
		"GBP",
		"JPY",
		"AUD",
		"KRW",
		"CAD",
		"NZD",
		"MOP",
		"HKD",
		"TWD",
		"AED",
		"SAR",
		"BRL",
		"CHF",
		"DKK",
		"SEK",
		"NOK",
		"IDR",
		"MYR",
		"PHP",
		"SGD",
		"THB",
		"INR",
		"RUB",
		"TRY",
		"ZAR",
	}
	for key, value := range data.Rates {
		if contains(cu_list, key) {
			continue
		}
		v, ok := globe.V4CruuencyMap.Load(key)
		if ok {
			existing, ok := v.(CurrencyInfo)
			if ok {
				existing.Rate = fmt.Sprint(value)
				existing.Update = "e" + getFormattedTime()
				globe.V4CruuencyMap.Store(key, existing)
			} else {
				fmt.Println("从er-api获取数据时类型断言失败:", key)
			}
		} else {
			fmt.Println("从er-api获取数据时发现新货币,未进一步处理:", key)
		}
	}
}

// 按每个国家代码获取汇率
func getDate(tcur string) (globe.Res, error) {
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
		return *res, err
	}

	if res.Success == "0" || res.Result.Tcur == "" {
		return *res, err
	}
	return *res, nil
}

// 把汇率发送到cloudflare中备份
func SendRateToCloudflare() {
	normalMap := make(map[string]interface{})
	globe.V4CruuencyMap.Range(func(key, value interface{}) bool {
		strKey, ok := key.(string) // 确保键的类型是 string
		if !ok {
			return true // 继续遍历
		}
		normalMap[strKey] = value
		return true // 继续遍历
	})

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(normalMap).
		Post("https://rate-back.oeoli.org/set_rate_v4")
	if err != nil {
		fmt.Printf("post_v4进CF备份时出错: %v\n", err)
		return
	}
	fmt.Printf("备份进CF的返回数据V4: %v\n", resp)

}

// 获取并格式化东八区的当前时间
// 获取并格式化东八区的当前时间
func getFormattedTime() string {
	// 创建东八区时区，UTC+8
	shanghaiZone := time.FixedZone("CST", 8*60*60)

	// 获取当前时间，并使用东八区时区
	now := time.Now().In(shanghaiZone)

	// 格式化时间
	const layout = "2006-01-02 15:04:05"
	return now.Format(layout)
}

// 定义一个汇率信息的结构体
type ExchangeInfo struct {
	BaseCode           string             `json:"base_oode"`
	Documentation      string             `json:"documentation"`
	Provider           string             `json:"provider"`
	Result             string             `json:"result"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeEolUnix        int64              `json:"time_eol_unix"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeLastUpdateUtc  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	TimeNextUpdateUtc  string             `json:"time_next_update_utc"`
	Rates              map[string]float64 `json:"rates"`
}

func getExchangData() (*ExchangeInfo, error) {
	url_ := "https://open.er-api.com/v6/latest/CNY" // 假设要GET请求的URL
	// proxyURL, err := url.Parse("http://localhost:1080")
	// if err != nil {
	// 	panic(err)
	// }

	// transport := &http.Transport{
	// 	Proxy: http.ProxyURL(proxyURL),
	// }

	// client := &http.Client{
	// 	Transport: transport,
	// }
	// 发起HTTP GET请求
	resp, err := http.Get(url_)
	if err != nil {
		// 处理错误
		fmt.Println("Error sending request to API endpoint. ", err)
		return nil, err
	}
	defer resp.Body.Close() // 确保在函数结束时关闭响应体

	// 用io.ReadAll读取响应体中的内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// 处理错误
		fmt.Println("Error reading response body. ", err)
		return nil, err
	}

	// 如果响应的内容类型是JSON，则可以进行解析
	if resp.Header.Get("Content-Type") == "application/json" {
		var data ExchangeInfo
		// 解析JSON数据到结构体
		if err := json.Unmarshal(body, &data); err != nil {
			// 处理错误
			fmt.Println("Error unmarshaling JSON to struct. ", err)
			return nil, err
		}
		if data.Result != "success" {
			fmt.Println("获取汇率出错")
			return nil, fmt.Errorf("获取汇率出错")
		}
		fmt.Println("open.er-api.com汇率更新时间:", data.TimeLastUpdateUtc)
		//把data.TimeLastUpdateUtc时间字符串，转变为本地时区字符串
		// utcTimestamp := int64(data.TimeLastUpdateUnix)
		// utcTime := time.Unix(utcTimestamp, 0)
		// localTime := utcTime.In(time.Local)
		// localTimeString := localTime.Format("2006-01-02 15:04:05")

		// TimeLastUpdate = localTimeString //这里把之前的时间更新掉
		// TimeLastGet = time.Now().Format("2006-01-02 15:04:05") //这里把之前的时间更新掉
		fmt.Println("汇率数据获取完成.")
		return &data, nil
	} else {
		// 如果响应不是JSON，则返回一个错误
		return nil, fmt.Errorf("response content-type is not application/json")
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
