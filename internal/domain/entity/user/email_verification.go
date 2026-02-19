package user

import "time"

const (
	VerificationTypeEmailVerification = "EMAIL_VERIFICATION"
	VerificationTypePasswordReset     = "PASSWORD_RESET"

	MaxOTPAttempts = 3
	OTPExpiry      = 10 * time.Minute
)

type EmailVerification struct {
	ID           int64     `gorm:"primaryKey;autoIncrement:false"`
	UserID       int64     `gorm:"not null;index"`
	Type         string    `gorm:"type:varchar(30);not null"`
	Token        string    `gorm:"uniqueIndex;type:varchar(255);not null"`
	OTPCode      string    `gorm:"type:varchar(6);not null"`
	Email        string    `gorm:"type:varchar(255);not null"`
	AttemptsUsed int       `gorm:"not null;default:0"`
	MaxAttempts  int       `gorm:"not null;default:3"`
	IsUsed       bool      `gorm:"not null;default:false"`
	ExpiresAt    time.Time `gorm:"not null"`
	CreatedAt    time.Time
}

func (EmailVerification) TableName() string {
	return "email_verifications"
}
