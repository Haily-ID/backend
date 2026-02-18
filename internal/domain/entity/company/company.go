package company

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	ID        int64          `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Code      string         `gorm:"uniqueIndex;type:varchar(50);not null" json:"code"`
	Address   string         `gorm:"type:text" json:"address"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Company) TableName() string {
	return "companies"
}
