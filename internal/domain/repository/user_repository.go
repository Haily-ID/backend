package repository

import (
	"context"

	"github.com/haily-id/engine/internal/domain/entity/user"
)

type UserRepository interface {
	Create(ctx context.Context, user *user.User) error
	FindByID(ctx context.Context, id int64) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	FindByIDWithCompanies(ctx context.Context, id int64) (*user.User, error)
	Update(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*user.User, error)
	AddCompany(ctx context.Context, userID, companyID int64) error
	RemoveCompany(ctx context.Context, userID, companyID int64) error
	GetCompaniesByUserID(ctx context.Context, userID int64) ([]*user.Company, error)
}
