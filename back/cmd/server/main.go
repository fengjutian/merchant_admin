package main

import (
	cmd "merchant_back/cmd"
	"merchant_back/internal/routes"
)

func main() {

	cmd.InitDB()

	// 使用路由配置
	r := routes.SetupRoutes()

	// 启动服务器，监听 8088 端口
	r.Run(":8088") // 默认监听 0.0.0.0:8080
}
