package services

import (
	"errors"
	"time"

	model "merchant_back/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	// 基础CRUD操作
	GetUserByID(id uint) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
	GetUsers(page, pageSize int) ([]*model.User, int64, error)

	// 认证相关方法
	Login(username, password string) (*model.User, error)
	Register(user *model.User) error
	ChangePassword(userID uint, oldPassword, newPassword string) error
	ResetPassword(email string, newPassword string) error
	UpdateLastLogin(userID uint) error

	// 业务方法
	GetUsersByRole(role string) ([]*model.User, error)
	GetUsersByStatus(status string) ([]*model.User, error)
	ActivateUser(id uint) error
	DeactivateUser(id uint) error
	SuspendUser(id uint) error

	// 初始化方法
	CreateDefaultAdmin() error
}

// userService 用户服务实现
type userService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := s.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *userService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建用户
func (s *userService) CreateUser(user *model.User) error {
	// 检查用户名是否已存在
	var existingUser model.User
	if err := s.db.Where("username = ? OR email = ?", user.Username, user.Email).First(&existingUser).Error; err == nil {
		return errors.New("用户名或邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 设置默认值
	if user.Role == "" {
		user.Role = model.UserRoleUser
	}
	if user.Status == "" {
		user.Status = model.UserStatusActive
	}

	return s.db.Create(user).Error
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(user *model.User) error {
	return s.db.Save(user).Error
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id uint) error {
	return s.db.Delete(&model.User{}, id).Error
}

// GetUsers 分页获取用户列表
func (s *userService) GetUsers(page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := s.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := s.db.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Login 用户登录
func (s *userService) Login(username, password string) (*model.User, error) {
	var user model.User

	// 支持用户名或邮箱登录
	err := s.db.Where("username = ? OR email = ?", username, username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 检查用户状态
	if !user.IsActive() {
		return nil, errors.New("用户账户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 更新最后登录时间
	if err := s.UpdateLastLogin(user.ID); err != nil {
		// 记录错误但不影响登录流程
		// 可以在这里添加日志记录
	}

	return &user, nil
}

// Register 用户注册
func (s *userService) Register(user *model.User) error {
	// 检查用户名长度
	if len(user.Username) < 3 || len(user.Username) > 50 {
		return errors.New("用户名长度必须在3-50个字符之间")
	}

	// 检查密码长度
	if len(user.Password) < 6 {
		return errors.New("密码长度不能少于6个字符")
	}

	// 检查邮箱格式
	if user.Email == "" {
		return errors.New("邮箱不能为空")
	}

	return s.CreateUser(user)
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 检查新密码长度
	if len(newPassword) < 6 {
		return errors.New("新密码长度不能少于6个字符")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = string(hashedPassword)
	return s.db.Save(user).Error
}

// ResetPassword 重置密码
func (s *userService) ResetPassword(email string, newPassword string) error {
	// 获取用户
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return err
	}

	// 检查新密码长度
	if len(newPassword) < 6 {
		return errors.New("新密码长度不能少于6个字符")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = string(hashedPassword)
	return s.db.Save(user).Error
}

// UpdateLastLogin 更新最后登录时间
func (s *userService) UpdateLastLogin(userID uint) error {
	now := time.Now()
	return s.db.Model(&model.User{}).Where("id = ?", userID).Update("last_login", now).Error
}

// GetUsersByRole 根据角色获取用户
func (s *userService) GetUsersByRole(role string) ([]*model.User, error) {
	var users []*model.User
	err := s.db.Where("role = ?", role).Find(&users).Error
	return users, err
}

// GetUsersByStatus 根据状态获取用户
func (s *userService) GetUsersByStatus(status string) ([]*model.User, error) {
	var users []*model.User
	err := s.db.Where("status = ?", status).Find(&users).Error
	return users, err
}

// ActivateUser 激活用户
func (s *userService) ActivateUser(id uint) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("status", model.UserStatusActive).Error
}

// DeactivateUser 禁用用户
func (s *userService) DeactivateUser(id uint) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("status", model.UserStatusInactive).Error
}

// SuspendUser 暂停用户
func (s *userService) SuspendUser(id uint) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("status", model.UserStatusSuspended).Error
}

// CreateDefaultAdmin 创建默认管理员账户
func (s *userService) CreateDefaultAdmin() error {
	// 检查是否已存在管理员
	var admin model.User
	err := s.db.Where("role = ?", model.UserRoleAdmin).First(&admin).Error
	if err == nil {
		return nil // 管理员已存在
	}

	// 创建默认管理员
	defaultAdmin := &model.User{
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  "admin123456", // 将被加密
		FirstName: "系统",
		LastName:  "管理员",
		Role:      model.UserRoleAdmin,
		Status:    model.UserStatusActive,
	}

	return s.CreateUser(defaultAdmin)
}
