package controllers

import (
	"net/http"

	"merchant_back/internal/common"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 这里写你的用户验证逻辑（例子先假装成功）
	if req.Username != "admin" || req.Password != "123456" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 假设 admin 的 userID 为 1
	token, _ := common.GenerateToken(1)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"type":  "Bearer",
	})
}
