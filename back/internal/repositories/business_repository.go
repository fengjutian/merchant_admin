package repositories

import (
	model "merchant_back/internal/models"

	"gorm.io/gorm"
)

// BusinessRepository 商家仓储接口
type BusinessRepository interface {
	// 基础CRUD操作
	Create(business *model.Business) error
	GetByID(id int) (*model.Business, error)
	GetAll() ([]*model.Business, error)
	Update(business *model.Business) error
	Delete(id int) error

	// 业务查询方法
	GetByEmail(email string) (*model.Business, error)
	GetByType(businessType string) ([]*model.Business, error)
	SearchByName(keyword string) ([]*model.Business, error)
	SearchByAddress(keyword string) ([]*model.Business, error)
	GetByRating(minRating float64) ([]*model.Business, error)
	GetByLocation(lat, lng float64, radius float64) ([]*model.Business, error)
	GetByStatus(status string) ([]*model.Business, error)

	// 分页查询
	GetWithPagination(page, pageSize int) ([]*model.Business, int64, error)

	// 批量操作
	BatchCreate(businesses []*model.Business) error
	BatchUpdate(businesses []*model.Business) error
	BatchDelete(ids []int) error

	// 统计方法
	Count() (int64, error)
	CountByType(businessType string) (int64, error)
	CountByStatus(status string) (int64, error)
}

// businessRepository 商家仓储实现
type businessRepository struct {
	db *gorm.DB
}

// NewBusinessRepository 创建商家仓储实例
func NewBusinessRepository(db *gorm.DB) BusinessRepository {
	return &businessRepository{
		db: db,
	}
}

// Create 创建商家
func (r *businessRepository) Create(business *model.Business) error {
	return r.db.Create(business).Error
}

// GetByID 根据ID获取商家
func (r *businessRepository) GetByID(id int) (*model.Business, error) {
	var business model.Business
	err := r.db.First(&business, id).Error
	if err != nil {
		return nil, err
	}
	return &business, nil
}

// GetAll 获取所有商家
func (r *businessRepository) GetAll() ([]*model.Business, error) {
	var businesses []*model.Business
	err := r.db.Where("status != ?", model.BusinessStatusSuspended).Find(&businesses).Error
	return businesses, err
}

// Update 更新商家
func (r *businessRepository) Update(business *model.Business) error {
	return r.db.Save(business).Error
}

// Delete 删除商家
func (r *businessRepository) Delete(id int) error {
	return r.db.Delete(&model.Business{}, id).Error
}

// GetByEmail 根据邮箱获取商家
func (r *businessRepository) GetByEmail(email string) (*model.Business, error) {
	var business model.Business
	err := r.db.Where("email = ?", email).First(&business).Error
	if err != nil {
		return nil, err
	}
	return &business, nil
}

// GetByType 根据类型获取商家列表
func (r *businessRepository) GetByType(businessType string) ([]*model.Business, error) {
	var business []*model.Business
	err := r.db.Where("type = ? AND status != ?", businessType, model.BusinessStatusSuspended).Find(&business).Error
	return business, err
}

// SearchByName 根据名称搜索商家
func (r *businessRepository) SearchByName(keyword string) ([]*model.Business, error) {
	var businesses []*model.Business
	err := r.db.Where("name LIKE ? AND status != ?", "%"+keyword+"%", model.BusinessStatusSuspended).Find(&businesses).Error
	return businesses, err
}

// SearchByAddress 根据地址搜索商家
func (r *businessRepository) SearchByAddress(keyword string) ([]*model.Business, error) {
	var businesses []*model.Business
	err := r.db.Where("address LIKE ? AND status != ?", "%"+keyword+"%", model.BusinessStatusSuspended).Find(&businesses).Error
	return businesses, err
}

// GetByRating 根据评分获取商家列表
func (r *businessRepository) GetByRating(minRating float64) ([]*model.Business, error) {
	var businesses []*model.Business
	err := r.db.Where("rating >= ? AND status != ?", minRating, model.BusinessStatusSuspended).Order("rating DESC").Find(&businesses).Error
	return businesses, err
}

// GetByLocation 根据位置获取附近商家（使用简化的距离计算）
func (r *businessRepository) GetByLocation(lat, lng float64, radius float64) ([]*model.Business, error) {
	var businesses []*model.Business

	// 简化的矩形范围查询（实际项目中可能需要更复杂的地理距离计算）
	latRange := radius / 111.0           // 大约1度纬度 = 111公里
	lngRange := radius / (111.0 * 0.866) // 粗略的经度范围计算

	query := r.db.Where("latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ? AND status != ?",
		-latRange, latRange, -lngRange, lngRange, model.BusinessStatusSuspended)

	err := query.Find(&businesses).Error
	return businesses, err
}

// GetByStatus 根据状态获取商家
func (r *businessRepository) GetByStatus(status string) ([]*model.Business, error) {
	var businesses []*model.Business
	err := r.db.Where("status = ?", status).Find(&businesses).Error
	return businesses, err
}

// GetWithPagination 分页获取商家列表
func (r *businessRepository) GetWithPagination(page, pageSize int) ([]*model.Business, int64, error) {
	var businesses []*model.Business
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.Model(&model.Business{}).Where("status != ?", model.BusinessStatusSuspended).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.Where("status != ?", model.BusinessStatusSuspended).Offset(offset).Limit(pageSize).Find(&businesses).Error
	return businesses, total, err
}

// BatchCreate 批量创建商家
func (r *businessRepository) BatchCreate(businesses []*model.Business) error {
	return r.db.CreateInBatches(businesses, 100).Error
}

// BatchUpdate 批量更新商家
func (r *businessRepository) BatchUpdate(businesses []*model.Business) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, business := range businesses {
			if err := tx.Save(business).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchDelete 批量删除商家
func (r *businessRepository) BatchDelete(ids []int) error {
	return r.db.Delete(&model.Business{}, ids).Error
}

// Count 获取商家总数
func (r *businessRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Business{}).Where("status != ?", model.BusinessStatusSuspended).Count(&count).Error
	return count, err
}

// CountByType 根据类型获取商家数量
func (r *businessRepository) CountByType(businessType string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Business{}).Where("type = ? AND status != ?", businessType, model.BusinessStatusSuspended).Count(&count).Error
	return count, err
}

// CountByStatus 根据状态获取商家数量
func (r *businessRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Business{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
