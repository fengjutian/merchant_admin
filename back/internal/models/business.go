package model

import "time"

// Business 商家模型
type Business struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`        // 商家名称
	Email       string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"` // 邮箱（唯一）
	Address     string    `gorm:"type:varchar(255);not null" json:"address"`     // 地址
	Type        string    `gorm:"type:varchar(100);not null" json:"type"`        // 商家类型
	Contact     string    `gorm:"type:varchar(255);not null" json:"contact"`     // 联系方式
	Rating      float64   `gorm:"default:0" json:"rating"`                       // 评分（0-5）
	Latitude    *float64  `gorm:"type:double" json:"latitude"`                   // 纬度（可空）
	Longitude   *float64  `gorm:"type:double" json:"longitude"`                  // 经度（可空）
	OtherInfo   *string   `gorm:"type:text" json:"otherInfo"`                    // 其他信息（可空）
	ImageBase64 *string   `gorm:"type:longtext" json:"imageBase64"`             // base64图片（可空）
	Description *string   `gorm:"type:text" json:"description"`                  // 描述（可空）
	Status      string    `gorm:"type:varchar(50);default:'active'" json:"status"` // 状态（active, inactive, suspended）
	Phone       string    `gorm:"type:varchar(20)" json:"phone"`                 // 电话号码
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`               // 创建时间
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`               // 更新时间
}

// TableName 指定表名
func (Business) TableName() string {
	return "businesses"
}

// BusinessStatus 商家状态常量
const (
	BusinessStatusActive    = "active"    // 活跃
	BusinessStatusInactive  = "inactive"  // 非活跃
	BusinessStatusSuspended = "suspended" // 暂停
)

// BusinessType 商家类型常量
const (
	BusinessTypeRestaurant = "restaurant" // 餐厅
	BusinessTypeRetail     = "retail"     // 零售
	BusinessTypeService    = "service"    // 服务
	BusinessTypeEntertainment = "entertainment" // 娱乐
	BusinessTypeOther      = "other"      // 其他
)
