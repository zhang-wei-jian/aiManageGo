package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	models "github.com/zhang-wei-jian/aiManageGo/modles"
	"net/http"
	"strconv"
)

// 中间件

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许跨域请求的来源
		//origin := c.Request.Header.Get("Origin")
		//c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		//c.Writer.Header().Set("Access-Control-Allow-Origin", "http://60.204.170.225:8503")
		// 设置允许跨域请求的来源为任何域
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 其他 CORS 头设置
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		} else {
			c.Next()
		}
	}
}

// 获取请求IP
func GetIp() gin.HandlerFunc {
	return func(c *gin.Context) {
		//gin 框架（包括目前业内所有的web框架），获取请求者IP有两种方式
		//
		//client_ip()
		//这种方式获取的是请求头中：X-Forwarded-For字段，以及client-ip字段，容易被伪造
		//request.remote_addr
		//这种方式为最真实的，但如果经过反向代理后，会显示为127.0.0.1

		ip := c.ClientIP()
		//ip := exnet.ClientPublicIP(c.Request)
		//ip := RemoteIP(c.Request)

		// 打印请求 IP
		fmt.Println("Request IP:", ip)

		//设置ip的次数
		go setRdbIpCount(ip)

		// 继续处理请求
		c.Next()
	}
}

func setRdbIpCount(ip string) {
	val := "0"
	val, err := models.Rdb.Get(ip).Result()
	if err != nil {
		fmt.Println("获取值错误:", err.Error())
		// 如果获取值错误，继续执行，并设置默认值为0
	}

	// 将字符串转换为整数
	count, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("转换值错误:", err.Error())
		// 如果转换值错误，将count设置为0
		count = 0
	}

	// 增加计数值
	count++

	// 将新值存储回 Redis
	err = models.Rdb.Set(ip, count, 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("获取IP值:", count)
}

func Response() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		//c.JSON(200, "今天使用额度用完了，明天再来看看吧")
		//c.JSON(200, gin.H{
		//	"model":  "chatglm3-6b",
		//	"id":     "",
		//	"object": "chat.completion",
		//	"choices": []interface{}{
		//		gin.H{
		//			"index": 0,
		//			"message": gin.H{
		//				"role":          "assistant",
		//				"content":       "我是一个基于 GPT-4.0 架构的语言模型，由 OpenAI 公司开发研发训练并获得许可。我的任务是针对用户的问题和要求提供适当的答复和支持。",
		//				"name":          nil,
		//				"function_call": nil,
		//			},
		//			"finish_reason": "stop",
		//		},
		//	},
		//	"created": 1708441994,
		//	"usage": gin.H{
		//		"prompt_tokens":     299,
		//		"total_tokens":      338,
		//		"completion_tokens": 39,
		//	},
		//})

		str := "不好意思今天不能用太多，明天再送你吧!See You"
		c.Header("Content-Type", "text/plain")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		for _, char := range str {
			c.String(http.StatusOK, string(char))
			c.Writer.Flush()

		}

	}
}
