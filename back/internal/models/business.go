package model

type business struct {
	ID          int      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string   `gorm:"type:varchar(255);not null" json:"name"`    // 业务名称
	Email       string   `gorm:"type:varchar(255);not null" json:"email"`   // 联系方式
	Address     string   `gorm:"type:varchar(255);not null" json:"address"` // 地址
	Type        string   `gorm:"type:varchar(100);not null" json:"type"`    // 类型
	Contact     string   `gorm:"type:varchar(255);not null" json:"contact"` // 通讯方式
	Rating      float64  `gorm:"default:0" json:"rating"`                   // 评价
	Latitude    *float64 `gorm:"type:double" json:"latitude"`               // 纬度（可空）
	Longitude   *float64 `gorm:"type:double" json:"longitude"`              // 经度（可空）
	OtherInfo   *string  `gorm:"type:text" json:"otherInfo"`                // 其他信息（可空）
	ImageBase64 *string  `gorm:"type:longtext" json:"imageBase64"`          // base64 图片（可空）
	Description *string  `gorm:"type:text" json:"description"`              // 描述（可空）
}
