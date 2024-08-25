package main

import (
	"fmt"
	"rate/api/globe"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func LoadV3Router(e *gin.Engine) {
	r := e.Group("/v3")
	r.GET("/rate", func(c *gin.Context) {
		cacheControlHeader := fmt.Sprintf("public, max-age=%d", viper.GetInt("cache.expire"))
		c.Header("Cache-Control", cacheControlHeader)
		c.Header("Data-Server", viper.GetString("DataServer"))
		normalMap := make(map[string]interface{})
		globe.V4CruuencyMap.Range(func(key, value interface{}) bool {
			strKey, ok := key.(string) // 确保键的类型是 string
			if !ok {
				return true // 继续遍历
			}
			normalMap[strKey] = value
			return true // 继续遍历
		})
		c.JSON(200, gin.H{"code": 20000, "data": normalMap, "message": "ok"})
	})
}
