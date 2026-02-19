package repository

import (
	"context"

	"github.com/haily-id/engine/internal/domain/entity/user"
)

type UserRepository interface {
	Create(ctx context.Context, u *user.User) error
	FindByID(ctx context.Context, id int64) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*user.User, error)
	Update(ctx context.Context, u *user.User) error
	Delete(ctx context.Context, id int64) error
}
