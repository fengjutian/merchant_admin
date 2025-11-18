package routes

import (
	"merchant_back/internal/controllers"
	"merchant_back/internal/middleware"
	"merchant_back/internal/repositories"
	"merchant_back/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 创建仓储层实例
	businessRepo := repositories.NewBusinessRepository(db)

	// 创建服务层实例
	businessService := services.NewBusinessService(businessRepo)

	// 创建控制器实例
	businessController := controllers.NewBusinessController(businessService)

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 商家路由
		businesses := api.Group("/business", middleware.AuthMiddleware())
		{
			businesses.GET("", businessController.GetBusinesses)            // 获取商家列表
			businesses.GET("/:id", businessController.GetBusiness)          // 获取单个商家
			businesses.POST("", businessController.CreateBusiness)          // 创建商家
			businesses.PUT("/:id", businessController.UpdateBusiness)       // 更新商家
			businesses.DELETE("/:id", businessController.DeleteBusiness)    // 删除商家
			businesses.GET("/search", businessController.SearchBusinesses)  // 搜索商家
			businesses.GET("/type", businessController.GetBusinessesByType) // 按类型获取商家
		}
	}

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

	return r
}
