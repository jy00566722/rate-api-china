package main

import (
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
)

type FeedBackInfo struct {
	Name    string `bson:"name" json:"name"`
	Url     string `bson:"url" json:"url"`
	Context string `bson:"context" json:"context"`
	Time    string `bson:"time" json:"time"`
}

func main() {
	r := gin.Default()
	r.GET("/good", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
		SendMail(FeedBackInfo{
			Name:    "报告员",
			Url:     "",
			Context: "Rate接口OK",
			Time:    time.Now().Format(time.DateTime),
		})
	})
	r.GET("/bad", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
		SendMail(FeedBackInfo{
			Name:    "报告员",
			Url:     "",
			Context: "Rate接口出错了!!!",
			Time:    time.Now().Format(time.DateTime),
		})
	})
	r.Run(":9999")
}

func SendMail(info FeedBackInfo) {
	e := email.NewEmail()
	e.From = "PengFeng <jy00566722@163.com>"
	e.To = []string{"jy00566723@gmail.com"}
	// e.Bcc = []string{"test_bcc@example.com"}
	// e.Cc = []string{"test_cc@example.com"}
	e.Subject = "汇率接口状态反馈"
	e.Text = []byte("Name:" + info.Name + "\n" + "URL:" + info.Url + "\n" + "内容:" + info.Context)
	// e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	err := e.Send("smtp.163.com:25", smtp.PlainAuth("", "jy00566722@163.com", "OTZJXGJJUQGMOBJR", "smtp.163.com"))
	if err != nil {
		fmt.Printf("发送邮件出错: %v\n", err)
	}
}
