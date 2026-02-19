package repository

import (
	"context"

	"github.com/haily-id/engine/internal/domain/entity/user"
)

type EmailVerificationRepository interface {
	Create(ctx context.Context, ev *user.EmailVerification) error
	FindByToken(ctx context.Context, token string) (*user.EmailVerification, error)
	FindActiveByUserIDAndType(ctx context.Context, userID int64, verType string) (*user.EmailVerification, error)
	MarkUsed(ctx context.Context, id int64) error
	IncrementAttempts(ctx context.Context, id int64) error
	InvalidateByUserIDAndType(ctx context.Context, userID int64, verType string) error
	DeleteByUserID(ctx context.Context, userID int64) error
}
