package user

import (
	"time"

	"gorm.io/gorm"
)

const (
	StatusPendingVerification = "PENDING_VERIFICATION"
	StatusActive              = "ACTIVE"
	StatusSuspended           = "SUSPENDED"

	GenderMale   = "MALE"
	GenderFemale = "FEMALE"
)

type User struct {
	ID              int64   `gorm:"primaryKey;autoIncrement:false"`
	Email           string  `gorm:"uniqueIndex;type:varchar(255);not null"`
	GoogleID        *string `gorm:"uniqueIndex;type:varchar(255)"`
	Password        *string `gorm:"type:varchar(255)"`
	Name            string  `gorm:"type:varchar(255);not null"`
	Phone           *string `gorm:"type:varchar(50)"`
	Gender          *string `gorm:"type:varchar(10)"`
	AvatarKey       *string `gorm:"type:varchar(500)"`
	Status          string  `gorm:"type:varchar(30);not null;default:'PENDING_VERIFICATION'"`
	EmailVerifiedAt *time.Time
	LastLoginAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
