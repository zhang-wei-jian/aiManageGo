package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	models "github.com/zhang-wei-jian/aiManageGo/modles"
	"net"
	"net/http"
	"strconv"
	"strings"
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

// RemoteIP 通过 RemoteAddr 获取 IP 地址， 只是一个快速解析方法。
func RemoteIP(r *http.Request) string {
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
