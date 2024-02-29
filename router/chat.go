package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	models "github.com/zhang-wei-jian/aiManageGo/modles"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 最多可用次數
var AvailableNumber = 10

// ip检测中间件。试用功能
func Chat(r *gin.Engine) {

	chatRouter := r.Group("/v1/chat")
	// 处理聊天请求/v1/chat/completions
	chatRouter.Use(ChatJWTMiddleware())
	chatRouter.POST("/completions", func(c *gin.Context) {

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

	c.Header("Content-Type", "text/plain")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	str := "不好意思今天不能用太多，明天再送你吧!See You"

	for _, char := range str {
		c.String(http.StatusOK, string(char))
		c.Writer.Flush()
	}
}

// 转发到响应服务器
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

// 获取和设置存ip到redis并且转发到代理服务器
func setRedisIpAndResponseProxyMessage(c *gin.Context) {
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
}

// 判断带不带token的聊天中间件
// 如果没问题直接return 结束当前中间件没有任何毛病
func ChatJWTMiddleware(allowPath ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//获取请求头token
		tokenHeader := ctx.GetHeader("Authorization")
		// 没登录,并且请求路径是聊天

		if (tokenHeader == "" || tokenHeader == "Bearer sk-xxxx") && ctx.Request.URL.Path == "/v1/chat/completions" {
			fmt.Println("no jwt token")
			//setRedisIpAndResponseProxyMessage(ctx)
			//ResponseProxyMessage(ctx)
			return
		} else {
			//	带请求头登录身份
			//segs := strings.SplitN(tokenHeader, " ", 2)
			segs := strings.Split(tokenHeader, " ")
			if len(segs) != 2 {
				// 格式不对，有人瞎搞
				fmt.Println("header - jwt - nil")
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			tokenStr := segs[1]
			//claims := &UserClaims{}
			claims := &UserClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(sign_key), nil
			})
			if err != nil {
				// token 不对，有人搞你
				fmt.Println("header - jwt - error")
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			if token == nil || !token.Valid || claims.Uid == 0 {
				// 按照道理来说，是不可能走到这一步的
				fmt.Println("header - jwt - aaaa")
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			now := time.Now()
			if claims.ExpiresAt.Time.Before(now) {
				// 过期了
				ctx.AbortWithStatus(http.StatusUnauthorized)
				fmt.Println("header - jwt - date")
				return
			}
			if claims.UserAgent != ctx.GetHeader("User-Agent") {
				// user agent 不相等
				ctx.AbortWithStatus(http.StatusUnauthorized)
				fmt.Println("header - jwt - aaaaac11a")
				fmt.Println(claims.UserAgent, "and", ctx.GetHeader("User-Agent"))

				return
			}
			// 为了演示，假设十秒钟刷新一次
			if claims.ExpiresAt.Time.Sub(now) < time.Second*50 {
				// 刷新
				claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
				tokenStr, err = token.SignedString([]byte(sign_key))
				if err != nil {
					// 因为刷新这个事情，并不是一定要做的，所以这里可以考虑打印日志
					// 暂时这样打印
					log.Println(err)
					return
				}
				ctx.Header("x-jwt-token", tokenStr)
				log.Println("刷新了 token")
			}

		}

	}

}
