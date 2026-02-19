package user

import (
	"strconv"

	userEntity "github.com/haily-id/engine/internal/domain/entity/user"
)

type UserDTO struct {
	ID              string  `json:"id"`
	Email           string  `json:"email"`
	Name            string  `json:"name"`
	Phone           *string `json:"phone"`
	Gender          *string `json:"gender"`
	AvatarKey       *string `json:"avatar_key"`
	Status          string  `json:"status"`
	EmailVerifiedAt *int64  `json:"email_verified_at"`
	LastLoginAt     *int64  `json:"last_login_at"`
	CreatedAt       int64   `json:"created_at"`
	UpdatedAt       int64   `json:"updated_at"`
}

func ToDTO(u *userEntity.User) UserDTO {
	dto := UserDTO{
		ID:        strconv.FormatInt(u.ID, 10),
		Email:     u.Email,
		Name:      u.Name,
		Phone:     u.Phone,
		Gender:    u.Gender,
		AvatarKey: u.AvatarKey,
		Status:    u.Status,
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	}

	if u.EmailVerifiedAt != nil {
		v := u.EmailVerifiedAt.Unix()
		dto.EmailVerifiedAt = &v
	}
	if u.LastLoginAt != nil {
		v := u.LastLoginAt.Unix()
		dto.LastLoginAt = &v
	}

	return dto
}
