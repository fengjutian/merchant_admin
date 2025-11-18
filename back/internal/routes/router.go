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
	userService := services.NewUserService(db)

	// 创建控制器实例
	businessController := controllers.NewBusinessController(businessService)
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(userService)

	// 认证路由
	r.POST("/login", authController.Login)
	r.POST("/api/auth/login", authController.Login)
	r.POST("/api/auth/refresh", authController.RefreshToken)
	r.POST("/api/auth/logout", authController.Logout)

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 商家路由
		businesses := api.Group("/business", middleware.AuthMiddleware())
		{
			businesses.GET("", businessController.GetBusinesses)                    // 获取商家列表
			businesses.GET("/:id", businessController.GetBusiness)                  // 获取单个商家
			businesses.POST("", businessController.CreateBusiness)                  // 创建商家
			businesses.PUT("/:id", businessController.UpdateBusiness)               // 更新商家
			businesses.DELETE("/:id", businessController.DeleteBusiness)             // 删除商家
			businesses.PUT("/:id/status", businessController.UpdateBusinessStatus)   // 更新商家状态
			businesses.GET("/status/:status", businessController.GetBusinessByStatus) // 根据状态获取商家
			businesses.GET("/search", businessController.SearchBusinesses)          // 搜索商家
			businesses.GET("/type", businessController.GetBusinessesByType)         // 按类型获取商家
			businesses.GET("/rating", businessController.GetBusinessesByRating)     // 按评分获取商家
			businesses.GET("/page", businessController.GetBusinessesWithPagination) // 分页获取商家
			businesses.GET("/nearby", businessController.GetNearbyBusinesses)       // 获取附近商家
			businesses.GET("/count", businessController.GetBusinessCount)           // 获取商家总数
			businesses.GET("/count/type", businessController.GetBusinessCountByType) // 根据类型获取商家数量
		}

		// 用户路由
		users := api.Group("/users")
		{
			users.GET("", userController.GetUsers)                        // 获取用户列表
			users.GET("/:id", userController.GetUser)                     // 获取单个用户
			users.POST("", userController.CreateUser)                     // 创建用户
			users.PUT("/:id", userController.UpdateUser)                  // 更新用户
			users.DELETE("/:id", userController.DeleteUser)               // 删除用户
			users.PUT("/:id/password", userController.ChangePassword)     // 修改密码
			users.PUT("/:id/status", userController.UpdateUserStatus)     // 更新用户状态
			users.GET("/role/:role", userController.GetUsersByRole)       // 根据角色获取用户列表
			users.GET("/status/:status", userController.GetUsersByStatus) // 根据状态获取用户列表
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
