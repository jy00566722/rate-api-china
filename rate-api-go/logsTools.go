package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 删除日志文件
func DelLogsFile() {
	fmt.Println("运行日志删除:")
	time.Sleep(10 * time.Second)
	reg, _ := regexp.Compile(`^(\d{4})-(\d{2})-(\d{2})-(\d{2})-(\d{2})-(\d{2})-gin.log`)
	// reg, _ := regexp.Compile(`^\d{4}-\d{2}-\d{2}-\d{2}-\d{d}-\d{2}-gin.log$`)
	files, err := os.ReadDir(LogsDir)
	if err != nil {
		log.Println(err)
	} else {
		for _, file := range files {
			fileName := file.Name()
			fmt.Printf("fileName: %v\n", fileName)
			if file.IsDir() {
				fmt.Printf("file.Name()这是目录: %v\n", fileName)
			} else {
				loc, err := time.LoadLocation("Asia/Shanghai")
				if err != nil {
					fmt.Println(err)
					return
				}
				strs := reg.FindStringSubmatch(fileName)
				now := time.Now()
				if len(strs) == 0 {
					// fmt.Printf("strs: %v\n", strs)
				} else {
					timeStr := strs[1] + `/` + strs[2] + `/` + strs[3] + " " + strs[4] + ":" + strs[5] + ":" + strs[6]
					fmt.Printf("timeStr: %v\n", timeStr)
					timeObj, err := time.ParseInLocation("2006/01/02 15:04:05", timeStr, loc)
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println(now.Sub(timeObj).Hours())
					if now.Sub(timeObj).Hours() > 0.001 {
						fmt.Printf("这个日志文件要删除了: %v\n", strs[0])
						err := os.Remove(LogsDir + "/" + strs[0])
						if err != nil {
							fmt.Println("删除失败!", err)
						}
					}
				}

			}
		}
	}
}

// GinLogger 接收gin框架默认的日志
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}
				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// 分析日志中的独立IP地址有多少个
func ShowIpNums() {
	file, err := os.Open("./logs/2022-07-26-20-05-52-gin.log")

	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	defer file.Close()
	allIp := make([]string, 0)
	fileScanner := bufio.NewScanner(file)
	reg, _ := regexp.Compile(`\d+\.\d+\.\d+\.\d+`)
	for fileScanner.Scan() {
		s := strings.Fields(fileScanner.Text())
		if len(s) >= 9 {
			aIp := s[9]
			res := reg.FindStringIndex(aIp)
			if res != nil {
				// fmt.Printf("aIp: %v\n", aIp)
				allIp = append(allIp, aIp)
			}
		}
	}
	result := RemoveDuplicateElement(allIp)
	fmt.Printf("result: %v\n", result)
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

}

func RemoveDuplicateElement(languages []string) []string {
	result := make([]string, 0, len(languages))
	temp := map[string]struct{}{}
	for _, item := range languages {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func RemoveDuplicatesInPlace(userIDs []string) []string {
	// 如果有0或1个元素，则返回切片本身。
	if len(userIDs) < 2 {
		return userIDs
	}

	//  使切片升序排序
	// sort.SliceStable(userIDs, func(i, j int) bool { return userIDs[i] < userIDs[j] })

	uniqPointer := 0

	for i := 1; i < len(userIDs); i++ {
		// 比较当前元素和唯一指针指向的元素
		//  如果它们不相同，则将项写入唯一指针的右侧。
		if userIDs[uniqPointer] != userIDs[i] {
			uniqPointer++
			userIDs[uniqPointer] = userIDs[i]
		}
	}

	return userIDs[:uniqPointer+1]
}

// router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
// 	// your custom format
// 	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
// 		param.ClientIP,  // 客户端IP
// 		param.TimeStamp.Format(time.RFC3339),  // 请求时间
// 		param.Method,  // 请求方法
// 		param.Path,  // 请求路径
// 		param.Request.Proto,  // 请求协议
// 		param.StatusCode,  // 请求状态码
// 		param.Latency,  // 请求时长
// 		param.Request.UserAgent(),  // 请求
// 		param.ErrorMessage,
// 	)
// }))

// func logsMid(param gin.LogFormatterParams) string {
// 	// your custom format
// 	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
// 		param.ClientIP,                       // 客户端IP
// 		param.TimeStamp.Format(time.RFC3339), // 请求时间
// 		param.Method,                         // 请求方法
// 		param.Path,                           // 请求路径
// 		param.Request.Proto,                  // 请求协议
// 		param.StatusCode,                     // 请求状态码
// 		param.Latency,                        // 请求时长
// 		param.Request.UserAgent(),            // 请求
// 		param.ErrorMessage,
// 	)
// }
