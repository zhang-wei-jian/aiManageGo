package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 中间件

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许跨域请求的来源
		//origin := c.Request.Header.Get("Origin")
		//c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://60.204.170.225:8503")
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
		// 获取请求 IP
		ip := c.Request.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = c.ClientIP()
		}
		// 打印请求 IP
		fmt.Println("Request IP:", ip)

		// 继续处理请求
		c.Next()
	}
}
