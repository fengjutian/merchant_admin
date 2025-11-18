package services

import (
	"errors"
	model "merchant_back/internal/models"
	"merchant_back/internal/repositories"
	"regexp"
	"strconv"
	"strings"
)

// BusinessService 商家服务接口
type BusinessService interface {
	// 基础CRUD操作
	GetBusinesses() ([]*model.Business, error)
	GetBusiness(id int) (*model.Business, error)
	CreateBusiness(business *model.Business) error
	UpdateBusiness(id int, business *model.Business) error
	DeleteBusiness(id int) error

	// 业务查询方法
	SearchBusinesses(query string) ([]*model.Business, error)
	GetBusinessesByType(businessType string) ([]*model.Business, error)
	GetBusinessByEmail(email string) (*model.Business, error)
	GetBusinessesByRating(minRating float64) ([]*model.Business, error)
	GetNearbyBusinesses(lat, lng, radius float64) ([]*model.Business, error)
	GetBusinessesByStatus(status string) ([]*model.Business, error)

	// 分页查询
	GetBusinessesWithPagination(page, pageSize int) ([]*model.Business, int64, error)

	// 批量操作
	BatchCreateBusinesses(businesses []*model.Business) error
	BatchUpdateBusinesses(businesses []*model.Business) error
	BatchDeleteBusinesses(ids []int) error

	// 统计方法
	GetBusinessCount() (int64, error)
	GetBusinessCountByType(businessType string) (int64, error)
	GetBusinessCountByStatus(status string) (int64, error)

	// 业务状态管理
	UpdateBusinessStatus(id int, status string) error
	ActivateBusiness(id int) error
	DeactivateBusiness(id int) error
	SuspendBusiness(id int) error
}

// businessService 商家服务实现
type businessService struct {
	businessRepo repositories.BusinessRepository
}

// NewBusinessService 创建商家服务实例
func NewBusinessService(businessRepo repositories.BusinessRepository) BusinessService {
	return &businessService{
		businessRepo: businessRepo,
	}
}

// validateEmail 验证邮箱格式
func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

// validatePhone 验证手机号格式（简单验证）
func validatePhone(phone string) bool {
	if phone == "" {
		return true // 手机号可以为空
	}
	pattern := `^1[3-9]\d{9}`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(phone)
}

// validateBusinessType 验证商家类型
func validateBusinessType(businessType string) bool {
	validTypes := []string{
		model.BusinessTypeRestaurant,
		model.BusinessTypeRetail,
		model.BusinessTypeService,
		model.BusinessTypeEntertainment,
		model.BusinessTypeOther,
	}

	for _, validType := range validTypes {
		if strings.EqualFold(businessType, validType) {
			return true
		}
	}
	return true // 允许自定义类型
}

// validateBusinessStatus 验证商家状态
func validateBusinessStatus(status string) bool {
	validStatuses := []string{
		model.BusinessStatusActive,
		model.BusinessStatusInactive,
		model.BusinessStatusSuspended,
	}

	for _, validStatus := range validStatuses {
		if strings.EqualFold(status, validStatus) {
			return true
		}
	}
	return false
}

// GetBusinesses 获取商家列表
func (s *businessService) GetBusinesses() ([]*model.Business, error) {
	return s.businessRepo.GetAll()
}

// GetBusiness 获取单个商家
func (s *businessService) GetBusiness(id int) (*model.Business, error) {
	if id <= 0 {
		return nil, errors.New("无效的商家ID")
	}

	business, err := s.businessRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("商家不存在")
	}

	return business, nil
}

// CreateBusiness 创建商家
func (s *businessService) CreateBusiness(business *model.Business) error {
	// 业务验证
	if business.Name == "" {
		return errors.New("商家名称不能为空")
	}
	if len(business.Name) > 255 {
		return errors.New("商家名称长度不能超过255个字符")
	}

	if business.Email == "" {
		return errors.New("邮箱不能为空")
	}
	if !validateEmail(business.Email) {
		return errors.New("邮箱格式不正确")
	}

	if business.Address == "" {
		return errors.New("地址不能为空")
	}
	if len(business.Address) > 255 {
		return errors.New("地址长度不能超过255个字符")
	}

	if business.Type == "" {
		return errors.New("商家类型不能为空")
	}
	if !validateBusinessType(business.Type) {
		return errors.New("商家类型无效")
	}

	if business.Contact == "" {
		return errors.New("联系方式不能为空")
	}
	if len(business.Contact) > 255 {
		return errors.New("联系方式长度不能超过255个字符")
	}

	if business.Phone != "" && !validatePhone(business.Phone) {
		return errors.New("手机号格式不正确")
	}

	if business.Rating < 0 || business.Rating > 5 {
		return errors.New("评分必须在0-5之间")
	}

	// 设置默认状态
	if business.Status == "" {
		business.Status = model.BusinessStatusActive
	} else if !validateBusinessStatus(business.Status) {
		return errors.New("商家状态无效")
	}

	// 检查邮箱是否已存在
	existingBusiness, err := s.businessRepo.GetByEmail(business.Email)
	if err == nil && existingBusiness != nil {
		return errors.New("邮箱已被使用")
	}

	// 设置默认评分
	if business.Rating == 0 {
		business.Rating = 0.0
	}

	return s.businessRepo.Create(business)
}

// UpdateBusiness 更新商家
func (s *businessService) UpdateBusiness(id int, business *model.Business) error {
	if id <= 0 {
		return errors.New("无效的商家ID")
	}

	// 检查商家是否存在
	existingBusiness, err := s.businessRepo.GetByID(id)
	if err != nil {
		return errors.New("商家不存在")
	}

	// 业务验证
	if business.Name == "" {
		return errors.New("商家名称不能为空")
	}
	if len(business.Name) > 255 {
		return errors.New("商家名称长度不能超过255个字符")
	}

	if business.Email == "" {
		return errors.New("邮箱不能为空")
	}
	if !validateEmail(business.Email) {
		return errors.New("邮箱格式不正确")
	}

	if business.Address == "" {
		return errors.New("地址不能为空")
	}
	if len(business.Address) > 255 {
		return errors.New("地址长度不能超过255个字符")
	}

	if business.Type == "" {
		return errors.New("商家类型不能为空")
	}
	if !validateBusinessType(business.Type) {
		return errors.New("商家类型无效")
	}

	if business.Contact == "" {
		return errors.New("联系方式不能为空")
	}
	if len(business.Contact) > 255 {
		return errors.New("联系方式长度不能超过255个字符")
	}

	if business.Phone != "" && !validatePhone(business.Phone) {
		return errors.New("手机号格式不正确")
	}

	if business.Rating < 0 || business.Rating > 5 {
		return errors.New("评分必须在0-5之间")
	}

	// 验证状态
	if business.Status != "" && !validateBusinessStatus(business.Status) {
		return errors.New("商家状态无效")
	}

	// 如果邮箱发生变化，检查新邮箱是否已被使用
	if business.Email != existingBusiness.Email {
		emailBusiness, err := s.businessRepo.GetByEmail(business.Email)
		if err == nil && emailBusiness != nil {
			return errors.New("邮箱已被使用")
		}
	}

	// 设置ID
	business.ID = id

	return s.businessRepo.Update(business)
}

// DeleteBusiness 删除商家
func (s *businessService) DeleteBusiness(id int) error {
	if id <= 0 {
		return errors.New("无效的商家ID")
	}

	// 检查商家是否存在
	_, err := s.businessRepo.GetByID(id)
	if err != nil {
		return errors.New("商家不存在")
	}

	return s.businessRepo.Delete(id)
}

// SearchBusinesses 搜索商家
func (s *businessService) SearchBusinesses(query string) ([]*model.Business, error) {
	if query == "" {
		return s.businessRepo.GetAll()
	}

	// 可以根据需求选择搜索策略
	// 这里简单实现：先按名称搜索，如果没有结果再按地址搜索
	businesses, err := s.businessRepo.SearchByName(query)
	if err != nil {
		return nil, err
	}

	// 如果名称搜索没有结果，尝试地址搜索
	if len(businesses) == 0 {
		businesses, err = s.businessRepo.SearchByAddress(query)
		if err != nil {
			return nil, err
		}
	}

	return businesses, nil
}

// GetBusinessesByType 根据类型获取商家
func (s *businessService) GetBusinessesByType(businessType string) ([]*model.Business, error) {
	if businessType == "" {
		return nil, errors.New("商家类型不能为空")
	}

	return s.businessRepo.GetByType(businessType)
}

// GetBusinessByEmail 根据邮箱获取商家
func (s *businessService) GetBusinessByEmail(email string) (*model.Business, error) {
	if email == "" {
		return nil, errors.New("邮箱不能为空")
	}

	return s.businessRepo.GetByEmail(email)
}

// GetBusinessesByRating 根据评分获取商家
func (s *businessService) GetBusinessesByRating(minRating float64) ([]*model.Business, error) {
	if minRating < 0 || minRating > 5 {
		return nil, errors.New("评分必须在0-5之间")
	}

	return s.businessRepo.GetByRating(minRating)
}

// GetNearbyBusinesses 获取附近商家
func (s *businessService) GetNearbyBusinesses(lat, lng, radius float64) ([]*model.Business, error) {
	if lat < -90 || lat > 90 {
		return nil, errors.New("纬度必须在-90到90之间")
	}
	if lng < -180 || lng > 180 {
		return nil, errors.New("经度必须在-180到180之间")
	}
	if radius <= 0 {
		return nil, errors.New("搜索半径必须大于0")
	}

	return s.businessRepo.GetByLocation(lat, lng, radius)
}

// GetBusinessesByStatus 根据状态获取商家
func (s *businessService) GetBusinessesByStatus(status string) ([]*model.Business, error) {
	if status == "" {
		return nil, errors.New("商家状态不能为空")
	}
	if !validateBusinessStatus(status) {
		return nil, errors.New("商家状态无效")
	}

	return s.businessRepo.GetByStatus(status)
}

// UpdateBusinessStatus 更新商家状态
func (s *businessService) UpdateBusinessStatus(id int, status string) error {
	if id <= 0 {
		return errors.New("无效的商家ID")
	}
	if status == "" {
		return errors.New("商家状态不能为空")
	}
	if !validateBusinessStatus(status) {
		return errors.New("商家状态无效")
	}

	// 获取现有商家信息
	business, err := s.businessRepo.GetByID(id)
	if err != nil {
		return errors.New("商家不存在")
	}

	// 更新状态
	business.Status = status
	return s.businessRepo.Update(business)
}

// GetBusinessesWithPagination 分页获取商家列表
func (s *businessService) GetBusinessesWithPagination(page, pageSize int) ([]*model.Business, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大页面大小
	}

	return s.businessRepo.GetWithPagination(page, pageSize)
}

// BatchCreateBusinesses 批量创建商家
func (s *businessService) BatchCreateBusinesses(businesses []*model.Business) error {
	if len(businesses) == 0 {
		return errors.New("商家列表不能为空")
	}

	// 验证每个商家
	for i, business := range businesses {
		if business.Name == "" {
			return errors.New("第" + strconv.Itoa(i+1) + "个商家名称不能为空")
		}
		if business.Email == "" {
			return errors.New("第" + strconv.Itoa(i+1) + "个商家邮箱不能为空")
		}
		if !validateEmail(business.Email) {
			return errors.New("第" + strconv.Itoa(i+1) + "个商家邮箱格式不正确")
		}
		// 检查邮箱是否重复
		existingBusiness, err := s.businessRepo.GetByEmail(business.Email)
		if err == nil && existingBusiness != nil {
			return errors.New("第" + strconv.Itoa(i+1) + "个商家邮箱已被使用")
		}
	}

	return s.businessRepo.BatchCreate(businesses)
}

// BatchUpdateBusinesses 批量更新商家
func (s *businessService) BatchUpdateBusinesses(businesses []*model.Business) error {
	if len(businesses) == 0 {
		return errors.New("商家列表不能为空")
	}

	// 验证每个商家
	for i, business := range businesses {
		if business.ID <= 0 {
			return errors.New("第" + strconv.Itoa(i+1) + "个商家ID无效")
		}
		if business.Name == "" {
			return errors.New("第" + strconv.Itoa(i+1) + "个商家名称不能为空")
		}
		if business.Email == "" {
			return errors.New("第" + strconv.Itoa(i+1) + "个商家邮箱不能为空")
		}
		if !validateEmail(business.Email) {
			return errors.New("第" + strconv.Itoa(i+1) + "个商家邮箱格式不正确")
		}

		// 检查商家是否存在
		_, err := s.businessRepo.GetByID(int(business.ID))
		if err != nil {
			return errors.New("第" + strconv.Itoa(i+1) + "个商家不存在")
		}
	}

	return s.businessRepo.BatchUpdate(businesses)
}

// BatchDeleteBusinesses 批量删除商家
func (s *businessService) BatchDeleteBusinesses(ids []int) error {
	if len(ids) == 0 {
		return errors.New("ID列表不能为空")
	}

	// 验证每个ID
	for _, id := range ids {
		if id <= 0 {
			return errors.New("包含无效的商家ID")
		}
		// 检查商家是否存在
		_, err := s.businessRepo.GetByID(id)
		if err != nil {
			return errors.New("ID为" + strconv.Itoa(id) + "的商家不存在")
		}
	}

	return s.businessRepo.BatchDelete(ids)
}

// GetBusinessCount 获取商家总数
func (s *businessService) GetBusinessCount() (int64, error) {
	return s.businessRepo.Count()
}

// GetBusinessCountByType 根据类型获取商家数量
func (s *businessService) GetBusinessCountByType(businessType string) (int64, error) {
	if businessType == "" {
		return 0, errors.New("商家类型不能为空")
	}

	return s.businessRepo.CountByType(businessType)
}

// GetBusinessCountByStatus 根据状态获取商家数量
func (s *businessService) GetBusinessCountByStatus(status string) (int64, error) {
	if status == "" {
		return 0, errors.New("商家状态不能为空")
	}
	if !validateBusinessStatus(status) {
		return 0, errors.New("商家状态无效")
	}

	return s.businessRepo.CountByStatus(status)
}

// ActivateBusiness 激活商家
func (s *businessService) ActivateBusiness(id int) error {
	business, err := s.GetBusiness(id)
	if err != nil {
		return err
	}

	business.Status = model.BusinessStatusActive
	return s.businessRepo.Update(business)
}

// DeactivateBusiness 停用商家
func (s *businessService) DeactivateBusiness(id int) error {
	business, err := s.GetBusiness(id)
	if err != nil {
		return err
	}

	business.Status = model.BusinessStatusInactive
	return s.businessRepo.Update(business)
}

// SuspendBusiness 暂停商家
func (s *businessService) SuspendBusiness(id int) error {
	business, err := s.GetBusiness(id)
	if err != nil {
		return err
	}

	business.Status = model.BusinessStatusSuspended
	return s.businessRepo.Update(business)
}
