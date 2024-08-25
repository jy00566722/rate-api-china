package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	//用gin框架，创建一个web服务器，监听8010端口,接口地址为v1/rate_usa
	//接口返回数据格式为json格式
	router := gin.Default()
	router.Use(cors.Default())
	// router.Use(printHeaders) //打印所有请求头
	router.Use(gin.Recovery())
	router.GET("/v1/rate_usa", func(c *gin.Context) {

		//返回数据
		c.JSON(200, gin.H{
			"code": 2000,
			"msg":  "success",
			"data": gin.H{
				"rate": 306221.32,
			},
		})
	})
	router.Run(":8010")
}

// printHeaders 中间件打印所有请求头
func printHeaders(c *gin.Context) {
	for key, value := range c.Request.Header {
		fmt.Println(key, ":", value)
	}
	// 之后调用路由处理函数
	c.Next()
}

//CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o btc-api-1219a
