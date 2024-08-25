package globe

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/go-resty/resty/v2"
)

// v4中保存汇率的sync.map
var V4CruuencyMap sync.Map

// v3中保存汇率的map
var CurrencyInfoMap CurrencyInfoMapType = make(CurrencyInfoMapType)

//v2中保存汇率的

var All_rate sync.Map
var All_rate_string string

var Mux sync.RWMutex // 读写锁 v2的

// 从map中获取数据转化为string,方便使用
func GetRateString() error {
	list, m := GetRateFromMap()
	if m > 130 {
		all_rate_byts, err := json.Marshal(list)
		if err != nil {
			fmt.Println("init,jsom.Marshal出错:", err)
			return err
		} else {
			fmt.Println("初始化时从map中获得的数据条数:", m)
			Mux.Lock()
			All_rate_string = string(all_rate_byts)
			Mux.Unlock()
			fmt.Println("更新内存中的数据all_rate_string成功!!")
		}
	} else {
		fmt.Println("初始化时从map中获得的数据条数-不足!!:", m, "没有保存到字符串中")
		return errors.New("初始化时从map中获得的数据条数-不足")
	}
	return nil
}

// 获取内存中的rate
func GetRateFromMap() (list []string, n int) {
	All_rate.Range(func(key, value any) bool {
		list = append(list, value.(string))
		n++
		return true
	})
	return
}

// 把all_rate_string发送到CF备份
func PostRateStringToCF(s string) (string, error) {
	// Create a Resty Client
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(s).
		Post("https://rate-back.oeoli.org/set_rate")
	if err != nil {
		fmt.Printf("post进CF备份时出错: %v\n", err)
		return "", err
	}
	fmt.Printf("resp: %v\n", resp)
	return string(resp.Body()), nil
}
