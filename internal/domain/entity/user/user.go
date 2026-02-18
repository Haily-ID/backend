package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64          `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Email     string         `gorm:"uniqueIndex;type:varchar(255);not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Role      string         `gorm:"type:varchar(50);not null;default:'user'" json:"role"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Many-to-many relationship with Company
	Companies []Company `gorm:"many2many:user_companies;" json:"companies,omitempty"`
}

type Company struct {
	ID        int64          `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Code      string         `gorm:"uniqueIndex;type:varchar(50);not null" json:"code"`
	Address   string         `gorm:"type:text" json:"address"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Many-to-many relationship with User
	Users []User `gorm:"many2many:user_companies;" json:"users,omitempty"`
}

type UserCompany struct {
	UserID    int64     `gorm:"primaryKey" json:"user_id"`
	CompanyID int64     `gorm:"primaryKey" json:"company_id"`
	JoinedAt  time.Time `gorm:"autoCreateTime" json:"joined_at"`
}

func (User) TableName() string {
	return "users"
}
