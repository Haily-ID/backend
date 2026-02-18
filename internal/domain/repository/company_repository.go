package repository

import (
	"context"

	"github.com/haily-id/engine/internal/domain/entity/company"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *company.Company) error
	FindByID(ctx context.Context, id int64) (*company.Company, error)
	FindByCode(ctx context.Context, code string) (*company.Company, error)
	Update(ctx context.Context, company *company.Company) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*company.Company, error)
}
