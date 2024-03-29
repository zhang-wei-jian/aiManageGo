package router

import (
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"net/url"
)

func Proxy(r *gin.Engine) {

	// 定义反向代理的目标地址
	//targetURL := "http://localhost:8898"
	//targetURL := "http://127.0.0.1:8000"
	targetURL := "http://60.204.170.225:8500"
	// 创建一个反向代理的 Transport
	target, _ := url.Parse(targetURL)
	proxy := httputil.NewSingleHostReverseProxy(target)

	// 定义路由，将所有请求都代理到目标地址
	r.Any("/*path", func(c *gin.Context) {

		//***********因为请求转发后还会cors坚决不跨域，如果设置多个就会导致报错，所以这里先删除掉*************
		// 删除默认的 CORS 头，确保不会转发后的代理产生冲突
		c.Writer.Header().Del("Access-Control-Allow-Origin")

		// 修改请求的 Host 头为目标地址的 Host，确保目标服务器能够正确识别请求
		c.Request.Host = target.Host
		// 执行反向代理
		proxy.ServeHTTP(c.Writer, c.Request)
	})

}
