package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type FeedBackInfo struct {
	Name    string `bson:"name" json:"name"`
	Url     string `bson:"url" json:"url"`
	Context string `bson:"context" json:"context"`
	Time    string `bson:"time" json:"time"`
}

func LoadFeedbackRouter(e *gin.Engine) {
	//接收feedback反馈信息
	e.POST("/v2/feedback", func(c *gin.Context) {
		var info FeedBackInfo
		c.ShouldBindJSON(&info)
		var ctx = context.Background()
		result, err := qClient.Database("rate").Collection("feedback").InsertOne(ctx, info)
		if err != nil {
			fmt.Printf("插入mongodb出错: %v\n", err)
		}

		go SendMail(info)
		c.JSON(200, gin.H{"code": 20000, "msg": "反馈成功", "data": result})
	})
}
