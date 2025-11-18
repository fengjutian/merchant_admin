package main

import (
	"log"
	"merchant_back/internal/config"
	"merchant_back/internal/controllers"
	"merchant_back/internal/routes"
	"merchant_back/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	InitDB()
	defer CloseDB()

	// 在 main.go 中
	userService := services.NewUserService(DB)
	controllers.NewAuthController(userService)

	// 使用路由配置
	r := routes.SetupRoutes(DB)

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
