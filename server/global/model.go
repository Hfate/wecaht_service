package global

import (
	"time"

	"gorm.io/gorm"
)

type BASEMODEL struct {
	ID        uint64         `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}
