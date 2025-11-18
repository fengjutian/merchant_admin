package repositories

import (
	model "merchant_back/internal/models"

	"gorm.io/gorm"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// 基础CRUD操作
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetAll() ([]*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error

	// 业务查询方法
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByUsernameOrEmail(usernameOrEmail string) (*model.User, error)
	GetByRole(role string) ([]*model.User, error)
	GetByStatus(status string) ([]*model.User, error)
	GetByPhone(phone string) (*model.User, error)

	// 搜索方法
	SearchByUsername(keyword string) ([]*model.User, error)
	SearchByEmail(keyword string) ([]*model.User, error)
	SearchByName(keyword string) ([]*model.User, error)

	// 分页查询
	GetWithPagination(page, pageSize int) ([]*model.User, int64, error)

	// 批量操作
	BatchCreate(users []*model.User) error
	BatchUpdate(users []*model.User) error
	BatchDelete(ids []uint) error

	// 统计方法
	Count() (int64, error)
	CountByRole(role string) (int64, error)
	CountByStatus(status string) (int64, error)

	// 业务方法
	UpdateLastLogin(userID uint) error
	UpdatePassword(userID uint, hashedPassword string) error
	UpdateStatus(userID uint, status string) error
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
}

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll 获取所有用户
func (r *userRepository) GetAll() ([]*model.User, error) {
	var users []*model.User
	err := r.db.Find(&users).Error
	return users, err
}

// Update 更新用户
func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsernameOrEmail 根据用户名或邮箱获取用户
func (r *userRepository) GetByUsernameOrEmail(usernameOrEmail string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByRole 根据角色获取用户列表
func (r *userRepository) GetByRole(role string) ([]*model.User, error) {
	var users []*model.User
	err := r.db.Where("role = ?", role).Find(&users).Error
	return users, err
}

// GetByStatus 根据状态获取用户列表
func (r *userRepository) GetByStatus(status string) ([]*model.User, error) {
	var users []*model.User
	err := r.db.Where("status = ?", status).Find(&users).Error
	return users, err
}

// GetByPhone 根据电话号码获取用户
func (r *userRepository) GetByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SearchByUsername 根据用户名搜索用户
func (r *userRepository) SearchByUsername(keyword string) ([]*model.User, error) {
	var users []*model.User
	err := r.db.Where("username LIKE ?", "%"+keyword+"%").Find(&users).Error
	return users, err
}

// SearchByEmail 根据邮箱搜索用户
func (r *userRepository) SearchByEmail(keyword string) ([]*model.User, error) {
	var users []*model.User
	err := r.db.Where("email LIKE ?", "%"+keyword+"%").Find(&users).Error
	return users, err
}

// SearchByName 根据姓名搜索用户
func (r *userRepository) SearchByName(keyword string) ([]*model.User, error) {
	var users []*model.User
	err := r.db.Where("first_name LIKE ? OR last_name LIKE ? OR CONCAT(first_name, ' ', last_name) LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Find(&users).Error
	return users, err
}

// GetWithPagination 分页获取用户列表
func (r *userRepository) GetWithPagination(page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}

// BatchCreate 批量创建用户
func (r *userRepository) BatchCreate(users []*model.User) error {
	return r.db.CreateInBatches(users, 100).Error
}

// BatchUpdate 批量更新用户
func (r *userRepository) BatchUpdate(users []*model.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, user := range users {
			if err := tx.Save(user).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchDelete 批量删除用户
func (r *userRepository) BatchDelete(ids []uint) error {
	return r.db.Delete(&model.User{}, ids).Error
}

// Count 获取用户总数
func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}

// CountByRole 根据角色获取用户数量
func (r *userRepository) CountByRole(role string) (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("role = ?", role).Count(&count).Error
	return count, err
}

// CountByStatus 根据状态获取用户数量
func (r *userRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// UpdateLastLogin 更新最后登录时间
func (r *userRepository) UpdateLastLogin(userID uint) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("last_login", gorm.Expr("NOW()")).Error
}

// UpdatePassword 更新用户密码
func (r *userRepository) UpdatePassword(userID uint, hashedPassword string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}

// UpdateStatus 更新用户状态
func (r *userRepository) UpdateStatus(userID uint, status string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("status", status).Error
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
