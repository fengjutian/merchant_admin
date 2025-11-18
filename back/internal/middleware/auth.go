package middleware

import (
	"net/http"
	"strings"

	"merchant_back/internal/common"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供Token"})
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims, err := common.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token无效或已过期"})
			c.Abort()
			return
		}

		// 保存 userID 到上下文
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
