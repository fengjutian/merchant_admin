package controllers

import (
	"net/http"
	"time"

	"merchant_back/internal/common"
	"merchant_back/internal/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthController 认证控制器
type AuthController struct {
	userService services.UserService
}

// NewAuthController 创建认证控制器实例
func NewAuthController(userService services.UserService) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string      `json:"token"`
	Type      string      `json:"type"`
	ExpiresAt int64       `json:"expiresAt"`
	User      interface{} `json:"user,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Login 用户登录
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest

	// 参数验证
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "参数验证失败",
			Error:   err.Error(),
		})
		return
	}

	// 查找用户
	user, err := ac.userService.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    401,
			Message: "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    401,
			Message: "用户名或密码错误",
		})
		return
	}

	// 检查用户状态
	if user.Status != "active" {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Code:    403,
			Message: "账户已被禁用",
		})
		return
	}

	// 生成 JWT token
	token, err := common.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "生成token失败",
			Error:   err.Error(),
		})
		return
	}

	// 更新最后登录时间
	user.UpdateLastLogin()
	if err := ac.userService.UpdateUser(user); err != nil {
		// 记录日志但不影响登录流程
		// log.Printf("更新用户最后登录时间失败: %v", err)
	}

	// 返回用户信息（不包含密码）
	userInfo := map[string]interface{}{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"role":      user.Role,
		"createdAt": user.CreatedAt,
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "登录成功",
		Data: LoginResponse{
			Token:     token,
			Type:      "Bearer",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			User:      userInfo,
		},
	})
}

// Logout 用户登出
func (ac *AuthController) Logout(c *gin.Context) {
	// 这里可以实现token黑名单机制
	// 或者记录登出日志

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "登出成功",
	})
}

// RefreshToken 刷新token
func (ac *AuthController) RefreshToken(c *gin.Context) {
	// 从header中获取token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    401,
			Message: "缺少认证token",
		})
		return
	}

	// 解析token
	tokenString := authHeader[7:] // 去掉 "Bearer "
	claims, err := common.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    401,
			Message: "token无效",
			Error:   err.Error(),
		})
		return
	}

	// 生成新token
	newToken, err := common.GenerateToken(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "刷新token失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "token刷新成功",
		Data: LoginResponse{
			Token:     newToken,
			Type:      "Bearer",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	})
}
