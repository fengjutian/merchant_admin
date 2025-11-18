package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建 Gin 路由引擎
	r := gin.Default()

	// 定义路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from Gin!",
			"status":  "success",
		})
	})

	// 定义另一个路由示例
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 启动服务器，监听 8080 端口
	r.Run(":8088") // 默认监听 0.0.0.0:8080
}
