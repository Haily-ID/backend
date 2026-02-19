package user

import (
	"context"
	"errors"
	"time"

	"github.com/haily-id/engine/internal/domain/entity/user"
	"github.com/haily-id/engine/internal/domain/repository"
	"gorm.io/gorm"
)

type emailVerificationRepository struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db *gorm.DB) repository.EmailVerificationRepository {
	return &emailVerificationRepository{db: db}
}

func (r *emailVerificationRepository) Create(ctx context.Context, ev *user.EmailVerification) error {
	return r.db.WithContext(ctx).Create(ev).Error
}

func (r *emailVerificationRepository) FindByToken(ctx context.Context, token string) (*user.EmailVerification, error) {
	var ev user.EmailVerification
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&ev).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("verification not found")
	}
	return &ev, err
}

func (r *emailVerificationRepository) FindActiveByUserIDAndType(ctx context.Context, userID int64, verType string) (*user.EmailVerification, error) {
	var ev user.EmailVerification
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ? AND is_used = false AND expires_at > ?", userID, verType, time.Now()).
		Order("created_at DESC").
		First(&ev).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("verification not found")
	}
	return &ev, err
}

func (r *emailVerificationRepository) MarkUsed(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&user.EmailVerification{}).
		Where("id = ?", id).
		Update("is_used", true).Error
}

func (r *emailVerificationRepository) IncrementAttempts(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&user.EmailVerification{}).
		Where("id = ?", id).
		UpdateColumn("attempts_used", gorm.Expr("attempts_used + 1")).Error
}

func (r *emailVerificationRepository) InvalidateByUserIDAndType(ctx context.Context, userID int64, verType string) error {
	return r.db.WithContext(ctx).
		Model(&user.EmailVerification{}).
		Where("user_id = ? AND type = ? AND is_used = false", userID, verType).
		Update("is_used", true).Error
}
