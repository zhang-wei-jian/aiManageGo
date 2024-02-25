package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	models "github.com/zhang-wei-jian/aiManageGo/modles"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

// 最多可用次數
var AvailableNumber = 10

func Chat(r *gin.Engine) {

	// 处理聊天请求
	r.POST("/v1/chat/completions", func(c *gin.Context) {

		//ResponseCustomMessage(c)
		//ResponseProxyMessage(c)
		//return
		ip := c.ClientIP()

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

		if count <= AvailableNumber {
			go func() {
				// 增加计数值
				count++
				// 将新值存储回 Redis
				err = models.Rdb.Set(ip, count, 0).Err()
				if err != nil {
					//panic(err)//程序退出
					fmt.Println("redis连接出现问题:", err)
				}

				fmt.Println("获取IP值:", count)
			}()
			//ai問答
			ResponseProxyMessage(c)

		} else {
			//不可以
			ResponseCustomMessage(c)
		}

	})

}

// 响应自定义内容
func ResponseCustomMessage(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	str := "不好意思今天不能用太多，明天再送你吧!See You"
	c.Header("Content-Type", "text/plain")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	for _, char := range str {
		c.String(http.StatusOK, string(char))
		c.Writer.Flush()
	}
}

func ResponseProxyMessage(c *gin.Context) {
	targetURL := "http://60.204.170.225:8500"
	// 创建一个反向代理的 Transport
	target, _ := url.Parse(targetURL)
	proxy := httputil.NewSingleHostReverseProxy(target)

	//***********因为请求转发后还会cors坚决不跨域，如果设置多个就会导致报错，所以这里先删除掉*************
	// 删除默认的 CORS 头，确保不会产生冲突

	c.Writer.Header().Del("Access-Control-Allow-Origin")

	// 修改请求的 Host 头为目标地址的 Host，确保目标服务器能够正确识别请求
	c.Request.Host = target.Host
	// 执行反向代理
	proxy.ServeHTTP(c.Writer, c.Request)

}
