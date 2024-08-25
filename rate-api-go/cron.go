package main

import (
	"fmt"
	"rate/api/v4rate"
	"time"

	"github.com/robfig/cron/v3"
)

// 定时任务
func MyCron() {
	fmt.Println("启动定时任务:....")
	c := cron.New(cron.WithSeconds())
	c.AddFunc("10 59 23 * * *", func() {
		fmt.Printf("运行定时任务-保存pop页面打开数据:%v\n", time.Now())
		CronSavePopTimesToDb()
	})
	if dataFrom == "google" {
		cron_google(c)
	} else if dataFrom == "now-api" {
		cron_nowapi(c)
	}

	c.AddFunc("20 5,35 * * * *", func() {
		v4rate.SendRateToCloudflare() //备份数据到cloudflare,每5分钟和35分钟执行一次
	})
	c.Start()
	fmt.Println("定时任务启动成功!....")
}

func cron_google(c *cron.Cron) {
	c.AddFunc("0 */10 * * * *", func() {
		v4rate.GetRateFromGoogleSheets() //每10分钟请求一次谷歌表格数据
	})
	c.AddFunc("30 10 */8 * * *", func() {
		v4rate.GetRateFromNowApi() //每8小时请求一次nowapi数据-不常见货币
	})
}
func cron_nowapi(c *cron.Cron) {
	c.AddFunc("10 5 8 * * *", func() {
		v4rate.GetRateFromErApi() //每天早上8点请求一次er-api数据
	})
	c.AddFunc("30 10 */12 * * *", func() {
		v4rate.GetRateFromNowApi() //每12小时请求一次nowapi数据-不常见货币
	})
	c.AddFunc("30 8 */1 * * *", func() {
		v4rate.GetRateFromNowApiCommon() //每1小时请求一次nowapi数据-常见货币
	})
}
