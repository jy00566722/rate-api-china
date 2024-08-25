package main

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func SendMail(info FeedBackInfo) {
	e := email.NewEmail()
	e.From = "PengFeng <jy00566722@163.com>"
	e.To = []string{"jy00566723@gmail.com"}
	// e.Bcc = []string{"test_bcc@example.com"}
	// e.Cc = []string{"test_cc@example.com"}
	e.Subject = "汇率转换新到反馈!"
	e.Text = []byte("Name:" + info.Name + "\n" + "URL:" + info.Url + "\n" + "内容:" + info.Context)
	// e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	err := e.Send("smtp.163.com:25", smtp.PlainAuth("", "jy00566722@163.com", "OTZJXGJJUQGMOBJR", "smtp.163.com"))
	if err != nil {
		fmt.Printf("发送邮件出错: %v\n", err)
	}
}
