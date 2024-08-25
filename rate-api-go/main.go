package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"rate/api/globe"
	"rate/api/v4rate"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/qiniu/qmgo"
	"github.com/spf13/viper"
)

var LogsDir string = "./logs" //日志目录
var Env string                //环境标示 dev or release
// var RateData map[string]string //汇率数据 定时获取

//go:embed templates/index.html
var indexHTML embed.FS
var PopPageOpenTimes atomic.Int32 //每天打开pop页面的总次数,每天定时写入数据库,再清零 原子操作不产生竞争
var PopPageIpTimes sync.Map       //每个ip打开的次数,每天保存
var qClient *qmgo.Client          //mongodb的客户端
var err error
var log_flag bool = false      //是否启用日志记录
var dataFrom string = "google" //数据来源 google 或是 now-api

func init() {
	GetConfig()
	// qClient, err = qmgo.NewClient(context.Background(), &qmgo.Config{Uri: "mongodb://rate:rate1234@t.deey.top:57890/?authSource=rate"})
	qClient, err = qmgo.NewClient(context.Background(), &qmgo.Config{Uri: viper.GetString("mongodbDns")})

	if err != nil {
		fmt.Printf("连接mongodb出错: %v\n", err)
	} else {
		fmt.Println("连接mongodb成功")
	}
	logsDirInfo, err := os.Stat(LogsDir)
	if err != nil {
		fmt.Println("目录不存在,需要新建")
		err := os.Mkdir(LogsDir, 0775)
		if err != nil {
			log.Panicf("新建日志目录失败")
		}
	} else {
		if !logsDirInfo.IsDir() {
			log.Panicf("存在%v,但不是目录", LogsDir)
		}
	}

	env := viper.GetString("env")
	// fmt.Printf("当前处于: %v模式\n", env)
	if env == "release" {
		Env = "release"
	} else {
		Env = "dev"
	}
	//获取数据来源,要判断是否为正确字符串，只能是google或是now-api
	dataFrom_ := viper.GetString("dataFrom")
	if dataFrom_ != "google" && dataFrom_ != "now-api" {
		fmt.Printf("数据来源只能是google或是now-api,现在使用默认的google\n")
	} else {
		dataFrom = dataFrom_
	}

	//获取缓存时间,判断是否能正确获取缓存时间
	cache_expire := viper.GetInt("cache.expire")
	fmt.Printf("缓存时间为: %v\n", cache_expire)
	if cache_expire < 20 {
		log.Panicf("不能发小于20秒的缓存时间,或者不能正确获取缓存时间")
	}
	//获取数据服务器名字,判断是否能正确获取
	data_server := viper.GetString("DataServer")
	fmt.Printf("服务器名字为: %v\n", data_server)
	if len(data_server) < 1 {
		log.Panicf("不能获取服务器名字")
	}
	//获取是否启用日志记录,因为nginx已经上了日志记录，在app中记录的意义不大
	log_flag = viper.GetBool("enableLogs")
	fmt.Printf("是否启用日志记录: %v\n", log_flag)
	//获取给客户端下发的请求频率
	periodInMinutes := viper.GetInt("periodInMinutes")
	fmt.Printf("客户端的请求频率为: %v 分钟\n", periodInMinutes)
	if periodInMinutes < 2 || periodInMinutes > 120 {
		log.Panicf("客户端的请求频率不正常")
	}
}
func main() {

	// 嵌入的文件作为slice of bytes读入
	content1, err := indexHTML.ReadFile("templates/index.html")
	if err != nil {
		// 请确保日志信息会输出到你能够看到的地方
		log.Fatalf("Failed to read embedded file: %v", err)
	}
	if dataFrom == "google" {
		v4rate.InitDataFromCFV4()        //从cloudflare获取数据
		v4rate.GetRateFromGoogleSheets() //从google sheets获取数据
		v4rate.GetRateFromNowApi()       //从nowapi获取数据-不常见货币
	} else if dataFrom == "now-api" {
		v4rate.InitDataFromCFV4()        //从cloudflare获取数据
		v4rate.GetRateFromErApi()        //从er-api获取数据
		v4rate.GetRateFromNowApi()       //从nowapi获取数据-不常见货币
		v4rate.GetRateFromNowApiCommon() //从nowapi获取数据-常见货币
	} else {
		log.Panicf("数据来源只能是google或是now-api,现在使用默认的google\n")
	}

	//日志文件的保存位置及分割时间
	f, _ := rotatelogs.New(
		"./logs/%Y%m%d%H%M.log",
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithMaxAge(time.Hour*24*10),
	)
	//按模式确定日志是保存还是直接输出
	env := viper.GetString("env")
	fmt.Printf("当前处于: %v模式\n", env)
	if env == "release" {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.MultiWriter(f)
	} else {
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}

	go MyCron() //定时任务 定时获取汇率API的数据

	r := gin.New()
	// pprof.Register(r)
	r.Use(cors.Default())
	if log_flag {
		r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			// 获取 Cloudflare 的 CF-Connecting-IP 头，如果它存在的话
			clientIP := param.Request.Header.Get("CF-Connecting-IP")
			if clientIP == "" {
				// 如果没有 CF-Connecting-IP，尝试获取 X-Forwarded-For 头
				clientIP = param.Request.Header.Get("X-Forwarded-For")
			}
			if clientIP == "" {
				// 如果没有 X-Forwarded-For，使用 Gin 的 ClientIP() 方法获取 IP
				clientIP = param.ClientIP
			}
			return fmt.Sprintf("[GINLOG]%s |%d |%10v |%s |%s |%s |%s |%s |%s\n",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"), // 请求时间
				param.StatusCode,          // 请求状态码
				param.Latency,             // 请求时长
				clientIP,                  // 客户端IP
				param.Method,              // 请求方法
				param.Request.Proto,       // 请求协议
				param.Path,                // 请求路径
				param.Request.UserAgent(), // 请求UA
				param.ErrorMessage,
			)
		}))
	} else {
		r.Use(gin.LoggerWithWriter(gin.DefaultErrorWriter))
	}

	r.Use(gin.Recovery())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello word")
	})

	LoadV1Router(r)       //注册V1版本中的路由处理
	LoadFeedbackRouter(r) //注册feedback的路由
	LoadErrorRouter(r)    //注册error的路由
	LoadV3Router(r)       //注册v3版本的路由

	r.GET("/custom-interface/call/Get_rate_redis", func(c *gin.Context) {

		var c_list []string
		globe.V4CruuencyMap.Range(func(key, value interface{}) bool {
			c, ok := value.(v4rate.CurrencyInfo)
			if !ok {
				fmt.Printf("value is not CurrencyInfoType")
				return true
			}
			parts := strings.Split(c.Name, " - ")
			if len(parts) != 2 {
				fmt.Println("数据格式错误:", c.Name)
				return true
			}
			ratenm := fmt.Sprintf("人民币/%s", parts[1])
			temp_c := globe.Result{
				Rate:   c.Rate,
				Ratenm: ratenm,
				Scur:   c.Scur,
				Status: c.Status,
				Tcur:   c.Tcur,
				Update: c.Update,
			}
			jsonData, err := json.Marshal(temp_c)
			if err != nil {
				fmt.Printf("json marshal error: %v\n", err)
				return true
			}
			c_list = append(c_list, string(jsonData))
			return true
		})

		cacheControlHeader := fmt.Sprintf("public, max-age=%d", viper.GetInt("cache.expire"))
		c.Header("Cache-Control", cacheControlHeader)
		c.Header("Data-Server", viper.GetString("DataServer"))
		c.JSON(http.StatusOK, c_list)

	})
	//发送所有的节点信息
	r.GET("/v2/all_nodes_and_date", func(c *gin.Context) {
		all := allNodes()
		cacheControlHeader := fmt.Sprintf("public, max-age=%d", viper.GetInt("cache.expire"))
		c.Header("Cache-Control", cacheControlHeader)
		c.Header("Data-Server", viper.GetString("DataServer"))
		c.JSON(200, gin.H{"code": 20000, "all_data": all})
	})
	//发送所有的节点信息
	r.GET("/v3/all_nodes_and_date", func(c *gin.Context) {
		all := allNodes()
		cacheControlHeader := fmt.Sprintf("public, max-age=%d", viper.GetInt("cache.expire"))
		c.Header("Cache-Control", cacheControlHeader)
		c.Header("Data-Server", viper.GetString("DataServer"))
		c.JSON(200, gin.H{"code": 20000, "all_data": all})
	})
	//接收pop页面打开的数据上报:次数 (后续提取更多内容)
	r.GET("/v2/feed_pop_page", func(c *gin.Context) {
		n := PopPageOpenTimes.Add(1)
		ip := c.ClientIP()
		if ip == "::1" {
			ip = "127.0.0.1"
		}
		//按IP记录，每个IP访问一次增加相应的次数
		v, ok := PopPageIpTimes.Load(ip)
		if ok {
			PopPageIpTimes.Store(ip, v.(int)+1)
		} else {
			PopPageIpTimes.Store(ip, 1)
		}
		c.JSON(200, gin.H{"code": 20000, "message": "ok", "data": n})
	})
	//返回pop页面打开次数的数据
	r.GET("/v2/show_pop_times", func(c *gin.Context) {
		m := make(map[string]int)
		PopPageIpTimes.Range(func(key, value any) bool {
			m[key.(string)] = value.(int)
			return true
		})
		c.JSON(200, gin.H{"code": 20000, "message": "ok", "data": PopPageOpenTimes.Load(), "dataIp": m})
	})
	//返回message信息，用于note页面及home页的提示,带main标志的是home页的提示,内容不能过长。
	r.GET("/v2/message", func(c *gin.Context) {
		message := viper.Get("message")
		c.JSON(200, gin.H{"code": 20000, "message": "ok", "data": message})
	})
	//卸载汇率插件时的接口,记录删除用户的信息 ip ua
	r.GET("/v2/uninstall", func(c *gin.Context) {
		type userInfo struct {
			Ip   string    `bson:"ip" json:"ip"`
			Ua   string    `bson:"ua" json:"ua"`
			Time time.Time `bson:"time" json:"time"`
			MyId string    `bson:"my_id" json:"my_id"`
		}
		clientIP := c.Request.Header.Get("CF-Connecting-IP")
		if clientIP == "" {
			// 如果没有 CF-Connecting-IP，尝试获取 X-Forwarded-For 头
			clientIP = c.Request.Header.Get("X-Forwarded-For")
		}
		if clientIP == "" {
			// 如果没有 X-Forwarded-For，使用 Gin 的 ClientIP() 方法获取 IP
			clientIP = c.ClientIP()
		}
		my_id := c.Query("my_id")
		var info = userInfo{
			Time: time.Now(),
			Ip:   clientIP,
			Ua:   c.Request.UserAgent(),
			MyId: my_id,
		}
		var ctx = context.Background()
		_, err := qClient.Database("rate").Collection("uninstall").InsertOne(ctx, info)
		if err != nil {
			fmt.Printf("插入mongodb出错: %v\n", err)
		}

		// c.String(http.StatusOK, "感谢你的使用,祝你生活愉快,如有意见及建议,请联系我:ideey88@gmail.com")
		c.Data(http.StatusOK, "text/html; charset=utf-8", content1)
	})
	//获取汇率--测试用-刚启动无汇率来源时用-CF备份也获取失败时用
	r.GET("/girl", func(c *gin.Context) {
		//CtronGetDate(true) //请求所有汇率- 随着v4接口的更新,这个接口已经废弃
		c.JSON(200, gin.H{"code": 20000, "message": "ok", "data": "执行请求所有汇率"})
	})
	r.GET("/pop", func(c *gin.Context) {
		CronSavePopTimesToDb() //保存pop页面打开次数
		c.JSON(200, gin.H{"code": 20000, "message": "ok", "data": "执行保存pop页面打开次数"})
	})
	//生成一个接口，用于让viper重新加载配置文件,使用ReadInConfig
	r.GET("/reload_config", func(c *gin.Context) {
		err := viper.ReadInConfig() // 查找并读取配置文件
		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{"code": 20000, "message": "error", "data": "重新加载配置文件失败"})
		}
		c.JSON(200, gin.H{"code": 20000, "message": "ok", "data": "重新加载配置文件OK"})
	})
	r.Run(":8001")
}

func allNodes() interface{} {

	var all_data = make(map[string]interface{})
	// all_data["aliexpress_element"] = viper.Get("aliexpress_element")
	all_data["aliexpress_nodes"] = viper.Get("aliexpress_nodes")
	// all_data["aliexpress_special_price"] = viper.Get("aliexpress_special_price")
	all_data["amazon_nodes"] = viper.Get("amazon_nodes")
	all_data["feedback_flag"] = viper.GetInt("feedback")
	all_data["lazada_main_price_nodes"] = viper.Get("lazada_main_price_nodes")
	all_data["lazada_nodes"] = viper.Get("lazada_nodes")
	all_data["ozon_nodes"] = viper.Get("ozon_nodes")
	all_data["ozon_nodeStr"] = viper.Get("ozon_nodeStr")
	all_data["qoo10jp_nodes"] = viper.Get("qoo10jp_nodes")
	all_data["shopee_nodes"] = viper.Get("shopee_nodes")
	all_data["gmarket_nodes"] = viper.Get("gmarket_nodes")
	all_data["temu_nodes"] = viper.Get("temu_nodes")
	all_data["temu_nodes_new"] = viper.Get("temu_nodes_new")
	all_data["wowmajp_nodes"] = viper.Get("wowmajp_nodes")
	all_data["rakuten_nodes"] = viper.Get("rakuten_nodes")
	all_data["vinted_nodes"] = viper.Get("vinted_nodes")
	all_data["miravia_nodes"] = viper.Get("miravia_nodes")
	all_data["kream_nodes"] = viper.Get("kream_nodes")
	all_data["mkd_nodes"] = viper.Get("mkd_nodes")
	all_data["wildberries"] = viper.Get("wildberries")
	all_data["coupang"] = viper.Get("coupang")
	all_data["walmart"] = viper.Get("walmart")
	all_data["message"] = viper.Get("message")
	all_data["periodInMinutes"] = viper.GetInt("periodInMinutes")
	return all_data
}

func GetConfig() {
	viper.SetConfigName("configMihu") // 配置文件名称(无扩展名)
	viper.AddConfigPath(".")          // 还可以在工作目录中查找配置
	viper.WatchConfig()
	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

			fmt.Println("no such config file")
		} else {

			fmt.Println("read config error")
		}
		fmt.Println(err)
	}

}

//CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rateApi
