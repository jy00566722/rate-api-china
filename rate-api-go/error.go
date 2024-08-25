package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

func LoadErrorRouter(e *gin.Engine) {
	e.POST("/v2/error", func(c *gin.Context) {

		var info ErrorInfo
		c.ShouldBindJSON(&info)

		//把收到的错误信息，写入mongodb
		var ctx = context.Background()
		_, err := qClient.Database("rate").Collection("error").InsertOne(ctx, info)
		if err != nil {
			fmt.Printf("插入mongodb出错_error: %v\n", err)
		}

		//给前端返回一个成功的json
		c.JSON(200, gin.H{"code": 20000, "msg": "反馈成功", "data": "ok"})

	})
}

type ErrorInfo struct {
	Ip        string `bson:"ip" json:"ip"`
	Ua        string `bson:"ua" json:"ua"`
	ErrorInfo string `bson:"error_info" json:"error_info"`
}
