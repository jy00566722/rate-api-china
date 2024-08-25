package v3rate

import (
	"embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rate/api/globe"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

type CurrencyInfo globe.CurrencyInfoType

type CurrencyInfoMap map[string]CurrencyInfo

var mux sync.Mutex

// 定义颜色代码
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
)

// 从json文件中初始化数据
func InitDateFromJson(girl3json embed.FS) {
	// 打开文件
	// 读取嵌入的 JSON 文件
	data, err := girl3json.ReadFile("templates/girl3.json")
	if err != nil {
		log.Fatal("这里出错了:", err)
	}

	// // 将文件内容读取到内存
	// jsonData, err := io.ReadAll(file)
	// if err != nil {
	// 	log.Fatal("Error when reading file: ", err)
	// } else {
	// 	log.Println("读取文件成功")
	// }

	// 解析到映射结构
	currencies := make(map[string]CurrencyInfo)
	err = json.Unmarshal(data, &currencies)
	if err != nil {
		log.Fatal("Error during unmarshal: ", err)
	}

	// 将映射结构转换为全局变量
	for code, info := range currencies {
		code := code
		info := info
		globe.CurrencyInfoMap[code] = globe.CurrencyInfoType{
			Name: info.Name,
			Rate: info.Rate,
			// RateNm:  info.RateNm,
			Scur:   info.Scur,
			Tcur:   info.Tcur,
			Update: info.Update,
			// FlagURL: info.FlagURL,
			Status: info.Status,
			Symbol: info.Symbol,
		}
	}
}

// 从cf中获取数据,更新全局变量
func UpdateCurrencyInfoMapFromCF() error {
	fmt.Println("从cf中获取数据,更新全局变量")
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
	for _, v := range rateStringList {
		v := v
		item_result := Result{}
		err = json.Unmarshal([]byte(v), &item_result)
		if err != nil {
			fmt.Printf("err2: %v\n", err)
			return err
		}
		name := item_result.Tcur + " - " + item_result.Ratenm

		v1, ok := globe.CurrencyInfoMap[item_result.Tcur]
		if ok {
			// mux.Lock()
			fmt.Printf("从cf取回数据中更新rate: %v,原汇率%v,新汇率%v\n", item_result.Tcur, v1.Rate, item_result.Rate)
			v1.Rate = item_result.Rate
			v1.Update = item_result.Update
			globe.CurrencyInfoMap[item_result.Tcur] = v1
			// mux.Unlock()

		} else {
			fmt.Println(Red, "从CF取回数据时发现新货币: ", name, Reset)
			// mux.Lock()
			globe.CurrencyInfoMap[item_result.Tcur] = globe.CurrencyInfoType{
				Name: name,
				Rate: item_result.Rate,
				// RateNm: item_result.Ratenm,
				Scur:   item_result.Scur,
				Tcur:   item_result.Tcur,
				Update: item_result.Update,
			}
			// mux.Unlock()
		}
	}
	return nil
}

// 从api接口中请求json，更新全局变量
func UpdateCurrencyInfoMapFromApi(result globe.Result) {
	time.Sleep(2 * time.Second)
	//判断globe.CurrencyInfoMap中是否已经存在该货币
	v, ok := globe.CurrencyInfoMap[result.Tcur]
	if ok {
		fmt.Printf("从v2接口发过来的数据中更新rate: %v,原汇率%v,新汇率%v\n", result.Tcur, v.Rate, result.Rate)
		// mux.Lock()
		v.Rate = result.Rate
		v.Update = result.Update
		globe.CurrencyInfoMap[result.Tcur] = v
		// mux.Unlock()
	} else {
		fmt.Println("从v2接口发过来的数据发现新货币: ", result.Tcur)
		// mux.Lock()
		globe.CurrencyInfoMap[result.Tcur] = globe.CurrencyInfoType{
			Name: result.Tcur + " - " + result.Ratenm,
			Rate: result.Rate,
			// RateNm: result.Ratenm,
			Scur:   result.Scur,
			Tcur:   result.Tcur,
			Update: result.Update,
		}
		// mux.Unlock()
	}
}

type Rate struct {
	CurrencyName string
	MiddleRate   string
	PubTime      string
}
type Rates map[string]Rate

// 从中国银行接口中获取数据,更新全局变量
func UpdateCurrencyInfoMapFromBOC() {

	fmt.Println("从中国银行接口中获取数据,更新全局变量")

	var allRates = make(Rates)
	url := "https://www.boc.cn/sourcedb/whpj/enindex_1619.html"
	// 获取 HTML 内容
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 查找具有指定宽度的表格
	doc.Find("table").Each(func(index int, tableHtml *goquery.Selection) {
		width, exists := tableHtml.Attr("width")
		if exists && width == "600" {
			// 在这个表格中提取汇率信息
			tableHtml.Find("tr").Each(func(indexTr int, rowHtml *goquery.Selection) {
				var rateData []string
				rowHtml.Find("td").Each(func(indexTd int, cellHtml *goquery.Selection) {
					text := strings.TrimSpace(cellHtml.Text())
					rateData = append(rateData, text)
				})

				if len(rateData) == 7 { // 确保是数据行
					currencyName := rateData[0]
					middleRate := rateData[5]
					pubTime := rateData[6]
					num, err := strconv.ParseFloat(middleRate, 64)
					if err != nil {
						// 错误处理
						fmt.Printf("err: %v\n", err)
						return
					}
					// 执行 1 除以这个数字，并乘以 100
					result := (1 / num) * 100

					// 将数字转换为字符串，保留小数点后八位
					resultStr := fmt.Sprintf("%.8f", result)
					// fmt.Printf("Currency: %s, Middle Rate: %s, Pub Time: %s\n", currencyName, middleRate, pubTime)
					//建立一个结构体，存储汇率数据
					rate := Rate{
						CurrencyName: currencyName,
						MiddleRate:   resultStr,
						PubTime:      pubTime,
					}
					//将汇率数据存储到map中
					allRates[currencyName] = rate
				}
			})
		}
	})
	// fmt.Printf("allRates: %v\n", allRates)
	// 先判断是否存在，存在则更新，不存在则添加到全局变量
	for _, v := range allRates {
		v := v
		//判断globe.CurrencyInfoMap中是否已经存在该货币
		v1, ok := globe.CurrencyInfoMap[v.CurrencyName]
		if ok {
			fmt.Printf("从中国银行取回数据中更新v3rate: %v,原汇率%v,新汇率%v\n", v.CurrencyName, v1.Rate, v.MiddleRate)
			v1.Rate = v.MiddleRate
			v1.Update = strings.Replace(v.PubTime, ".", "-", -1)
			v1.Hot = 1
			globe.CurrencyInfoMap[v.CurrencyName] = v1

		} else {
			fmt.Println("从中国银行取回数据时发现新货币: ", v.CurrencyName)
			globe.CurrencyInfoMap[v.CurrencyName] = globe.CurrencyInfoType{
				Name: v.CurrencyName,
				Rate: v.MiddleRate,
				// RateNm: v.CurrencyName,
				Scur:   "CNY",
				Tcur:   v.CurrencyName,
				Update: strings.Replace(v.PubTime, ".", "-", -1),
				Hot:    1,
			}
		}
	}

	//更新v2接口中的数据,先更新map，再更新字符串
	UpdateRateInfoMapFromApiv2(allRates)

}

// 从google sheets中获取汇率
func GetRateFromGoogleSheets() {
	// 发送 HTTP GET 请求
	// 兑人民币的汇率连接
	// https://docs.google.com/spreadsheets/d/e/2PACX-1vQS4NljENMDBH7xGLQ7OpmogFqm7mnwMg_W6MctPGxfaC0GJbitgFejoRUCAlmo9fh1k75DfXhOeMRN/pub?gid=0&single=true&output=csv
	// 兑美元的汇率连接
	// https://docs.google.com/spreadsheets/d/e/2PACX-1vQlmuo1VKRfqHEmzRslTHf1xxAukaMo8TK3kNM7FUGqAhfjH16RoqUb0gJOf71oOzEfaJhQVx6SsyON/pub?gid=0&single=true&output=csv
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
	// ExchangeRate 结构体表示汇率数据
	type ExchangeRate struct {
		Rate float64
		Code string
	}
	// 解析 CSV 数据到结构体
	var exchangeRates []ExchangeRate
	for _, record := range records {
		// 解析汇率
		rate, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			fmt.Println("Error parsing rate:", err)
			continue
		}

		// 创建 ExchangeRate 结构体
		exchangeRate := ExchangeRate{
			Rate: rate,
			Code: record[1],
		}

		// 添加到切片
		exchangeRates = append(exchangeRates, exchangeRate)
	}

	// 更新汇率数据
	for _, er := range exchangeRates {
		// fmt.Printf("Symbol: %s, Rate: %f\n", er.Code, er.Rate)
		v2, ok := globe.CurrencyInfoMap[er.Code]
		if ok {
			v2.Rate = fmt.Sprintf("%f", er.Rate)
			v2.Update = time.Now().Format("2006-01-02 15:04:05")
			globe.CurrencyInfoMap[er.Code] = v2
			fmt.Printf("从google sheets取回数据中更新v3rate: %v,原汇率%v,新汇率%v\n", er.Code, v2.Rate, er.Rate)
		} else {
			fmt.Println("从google sheets取回数据时发现新货币: ", er.Code)
			globe.CurrencyInfoMap[er.Code] = globe.CurrencyInfoType{
				Name:   er.Code,
				Rate:   fmt.Sprintf("%f", er.Rate),
				Scur:   "CNY",
				Tcur:   er.Code,
				Update: time.Now().Format("2006-01-02 15:04:05"),
				Hot:    1,
			}
		}

	}
}
