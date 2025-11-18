package controllers

import (
	"net/http"
	"strconv"

	model "merchant_back/internal/models"
	"merchant_back/internal/services"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService services.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取所有用户信息，支持分页参数
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} SuccessResponse{data=map[string]interface{}} "获取成功"
// @Failure 500 {object} ErrorResponse "获取失败"
// @Router /api/v1/users [get]
func (uc *UserController) GetUsers(c *gin.Context) {
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取用户列表
	users, total, err := uc.userService.GetUsers(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "获取用户列表失败",
			Error:   err.Error(),
		})
		return
	}

	// 转换用户数据（不包含密码）
	userList := make([]map[string]interface{}, len(users))
	for i, user := range users {
		userList[i] = map[string]interface{}{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"phone":     user.Phone,
			"avatar":    user.Avatar,
			"role":      user.Role,
			"status":    user.Status,
			"lastLogin": user.LastLogin,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "获取用户列表成功",
		Data: map[string]interface{}{
			"users": userList,
			"pagination": map[string]interface{}{
				"page":     page,
				"pageSize": pageSize,
				"total":    total,
				"pages":    (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GetUser 获取单个用户信息
// @Summary 获取用户详情
// @Description 根据用户ID获取单个用户的详细信息
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} SuccessResponse{data=map[string]interface{}} "获取成功"
// @Failure 400 {object} ErrorResponse "无效的用户ID"
// @Failure 404 {object} ErrorResponse "用户不存在"
// @Router /api/v1/users/{id} [get]
func (uc *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "无效的用户ID",
		})
		return
	}

	user, err := uc.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    404,
			Message: "用户不存在",
			Error:   err.Error(),
		})
		return
	}

	// 返回用户信息（不包含密码）
	userInfo := map[string]interface{}{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"phone":     user.Phone,
		"avatar":    user.Avatar,
		"role":      user.Role,
		"status":    user.Status,
		"lastLogin": user.LastLogin,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "获取用户信息成功",
		Data:    userInfo,
	})
}

// CreateUser 创建用户
// @Summary 创建新用户
// @Description 创建一个新的用户账户
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.UserRegisterRequest true "用户信息"
// @Success 201 {object} SuccessResponse{data=map[string]interface{}} "创建成功"
// @Failure 400 {object} ErrorResponse "参数验证失败或创建失败"
// @Router /api/v1/users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	var req model.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "参数验证失败",
			Error:   err.Error(),
		})
		return
	}

	// 创建用户对象
	user := &model.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	}

	// 创建用户
	if err := uc.userService.Register(user); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "创建用户失败",
			Error:   err.Error(),
		})
		return
	}

	// 返回创建的用户信息（不包含密码）
	userInfo := map[string]interface{}{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"phone":     user.Phone,
		"role":      user.Role,
		"status":    user.Status,
		"createdAt": user.CreatedAt,
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Code:    201,
		Message: "创建用户成功",
		Data:    userInfo,
	})
}

// UpdateUser 更新用户信息
// @Summary 更新用户信息
// @Description 根据用户ID更新用户的基本信息
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param user body models.UserUpdateRequest true "用户更新信息"
// @Success 200 {object} SuccessResponse{data=map[string]interface{}} "更新成功"
// @Failure 400 {object} ErrorResponse "无效的用户ID或参数验证失败"
// @Failure 404 {object} ErrorResponse "用户不存在"
// @Failure 500 {object} ErrorResponse "更新失败"
// @Router /api/v1/users/{id} [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "无效的用户ID",
		})
		return
	}

	var req model.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "参数验证失败",
			Error:   err.Error(),
		})
		return
	}

	// 获取现有用户
	user, err := uc.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    404,
			Message: "用户不存在",
			Error:   err.Error(),
		})
		return
	}

	// 更新用户信息
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone
	user.Avatar = req.Avatar

	// 保存更新
	if err := uc.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "更新用户失败",
			Error:   err.Error(),
		})
		return
	}

	// 返回更新后的用户信息
	userInfo := map[string]interface{}{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"phone":     user.Phone,
		"avatar":    user.Avatar,
		"role":      user.Role,
		"status":    user.Status,
		"lastLogin": user.LastLogin,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "更新用户成功",
		Data:    userInfo,
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 根据用户ID删除用户账户
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} SuccessResponse "删除成功"
// @Failure 400 {object} ErrorResponse "无效的用户ID"
// @Failure 404 {object} ErrorResponse "用户不存在"
// @Failure 500 {object} ErrorResponse "删除失败"
// @Router /api/v1/users/{id} [delete]
func (uc *UserController) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "无效的用户ID",
		})
		return
	}

	// 检查用户是否存在
	_, err = uc.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    404,
			Message: "用户不存在",
			Error:   err.Error(),
		})
		return
	}

	// 删除用户
	if err := uc.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "删除用户失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "删除用户成功",
	})
}

// ChangePassword 修改密码
// @Summary 修改用户密码
// @Description 根据用户ID修改用户密码
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param password body models.PasswordChangeRequest true "密码信息"
// @Success 200 {object} SuccessResponse "修改成功"
// @Failure 400 {object} ErrorResponse "无效的用户ID或参数验证失败"
// @Router /api/v1/users/{id}/password [put]
func (uc *UserController) ChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "无效的用户ID",
		})
		return
	}

	var req model.PasswordChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "参数验证失败",
			Error:   err.Error(),
		})
		return
	}

	// 修改密码
	if err := uc.userService.ChangePassword(uint(id), req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "修改密码失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "修改密码成功",
	})
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Description 根据用户ID更新用户状态（active/inactive/suspended）
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param status body object{status=string} true "用户状态"
// @Success 200 {object} SuccessResponse "更新成功"
// @Failure 400 {object} ErrorResponse "无效的用户ID或状态值"
// @Failure 500 {object} ErrorResponse "更新失败"
// @Router /api/v1/users/{id}/status [put]
func (uc *UserController) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "无效的用户ID",
		})
		return
	}

	type StatusRequest struct {
		Status string `json:"status" binding:"required,oneof=active inactive suspended"`
	}

	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "参数验证失败",
			Error:   err.Error(),
		})
		return
	}

	switch req.Status {
	case model.UserStatusActive:
		err = uc.userService.ActivateUser(uint(id))
	case model.UserStatusInactive:
		err = uc.userService.DeactivateUser(uint(id))
	case model.UserStatusSuspended:
		err = uc.userService.SuspendUser(uint(id))
	default:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "无效的状态值",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "更新用户状态失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "更新用户状态成功",
	})
}

// GetUsersByRole 根据角色获取用户列表
// @Summary 根据角色获取用户列表
// @Description 根据用户角色获取用户列表
// @Tags users
// @Accept json
// @Produce json
// @Param role path string true "用户角色" Enums(admin,user,merchant)
// @Success 200 {object} SuccessResponse{data=[]map[string]interface{}} "获取成功"
// @Failure 400 {object} ErrorResponse "角色参数不能为空"
// @Failure 500 {object} ErrorResponse "获取失败"
// @Router /api/v1/users/role/{role} [get]
func (uc *UserController) GetUsersByRole(c *gin.Context) {
	role := c.Param("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "角色参数不能为空",
		})
		return
	}

	users, err := uc.userService.GetUsersByRole(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "获取用户列表失败",
			Error:   err.Error(),
		})
		return
	}

	// 转换用户数据
	userList := make([]map[string]interface{}, len(users))
	for i, user := range users {
		userList[i] = map[string]interface{}{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"role":      user.Role,
			"status":    user.Status,
			"createdAt": user.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "获取用户列表成功",
		Data:    userList,
	})
}

// GetUsersByStatus 根据状态获取用户列表
// @Summary 根据状态获取用户列表
// @Description 根据用户状态获取用户列表
// @Tags users
// @Accept json
// @Produce json
// @Param status path string true "用户状态" Enums(active,inactive,suspended)
// @Success 200 {object} SuccessResponse{data=[]map[string]interface{}} "获取成功"
// @Failure 400 {object} ErrorResponse "状态参数不能为空"
// @Failure 500 {object} ErrorResponse "获取失败"
// @Router /api/v1/users/status/{status} [get]
func (uc *UserController) GetUsersByStatus(c *gin.Context) {
	status := c.Param("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "状态参数不能为空",
		})
		return
	}

	users, err := uc.userService.GetUsersByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "获取用户列表失败",
			Error:   err.Error(),
		})
		return
	}

	// 转换用户数据
	userList := make([]map[string]interface{}, len(users))
	for i, user := range users {
		userList[i] = map[string]interface{}{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"role":      user.Role,
			"status":    user.Status,
			"createdAt": user.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "获取用户列表成功",
		Data:    userList,
	})
}
