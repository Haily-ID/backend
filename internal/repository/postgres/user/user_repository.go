package user

import (
	"context"
	"errors"

	"github.com/haily-id/engine/internal/domain/entity/user"
	"github.com/haily-id/engine/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &u, err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &u, err
}

func (r *userRepository) FindByGoogleID(ctx context.Context, googleID string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("google_id = ?", googleID).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &u, err
}

func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&user.User{}, id).Error
}

