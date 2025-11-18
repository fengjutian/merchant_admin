package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"username"` // 用户名（唯一）
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`   // 邮箱（唯一）
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`                    // 密码（不返回给前端）
	FirstName string    `gorm:"type:varchar(100)" json:"firstName"`                      // 名
	LastName  string    `gorm:"type:varchar(100)" json:"lastName"`                       // 姓
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`                           // 电话号码
	Avatar    string    `gorm:"type:varchar(500)" json:"avatar"`                        // 头像URL
	Role      string    `gorm:"type:varchar(50);default:'user'" json:"role"`            // 角色（admin, user, merchant）
	Status    string    `gorm:"type:varchar(50);default:'active'" json:"status"`        // 状态（active, inactive, suspended）
	LastLogin *time.Time `gorm:"type:datetime" json:"lastLogin"`                        // 最后登录时间
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`                        // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`                        // 更新时间
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserStatus 用户状态常量
const (
	UserStatusActive    = "active"    // 活跃
	UserStatusInactive  = "inactive"  // 非活跃
	UserStatusSuspended = "suspended" // 暂停
)

// UserRole 用户角色常量
const (
	UserRoleAdmin    = "admin"    // 管理员
	UserRoleUser     = "user"     // 普通用户
	UserRoleMerchant = "merchant" // 商家
)

// GetFullName 获取用户全名
func (u *User) GetFullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	if u.LastName != "" {
		return u.LastName
	}
	return u.Username
}

// IsActive 检查用户是否活跃
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsAdmin 检查用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// IsMerchant 检查用户是否为商家
func (u *User) IsMerchant() bool {
	return u.Role == UserRoleMerchant
}

// UpdateLastLogin 更新最后登录时间
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名或邮箱
	Password string `json:"password" binding:"required"` // 密码
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`  // 用户名
	Email     string `json:"email" binding:"required,email"`           // 邮箱
	Password  string `json:"password" binding:"required,min=6"`        // 密码
	FirstName string `json:"firstName" binding:"max=100"`              // 名
	LastName  string `json:"lastName" binding:"max=100"`               // 姓
	Phone     string `json:"phone" binding:"max=20"`                   // 电话号码
}

// UserUpdateRequest 用户更新请求
type UserUpdateRequest struct {
	FirstName string `json:"firstName" binding:"max=100"` // 名
	LastName  string `json:"lastName" binding:"max=100"`  // 姓
	Phone     string `json:"phone" binding:"max=20"`      // 电话号码
	Avatar    string `json:"avatar" binding:"max=500"`   // 头像URL
}

// PasswordChangeRequest 密码修改请求
type PasswordChangeRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"` // 旧密码
	NewPassword string `json:"newPassword" binding:"required,min=6"` // 新密码
}

// PasswordResetRequest 密码重置请求
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"` // 邮箱
}