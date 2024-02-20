package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zhang-wei-jian/aiManageGo/router"
)

func main() {
	// 创建一个新的 Gin 实例
	r := gin.Default()

	//获取IP存入redis
	r.Use(router.GetIp())

	//转发请求
	router.Proxy(r)

	// 启动 Gin 服务器，监听端口
	r.Run(":8080")
}
