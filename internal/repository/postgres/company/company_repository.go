package company

import (
	"context"
	"errors"

	"github.com/haily-id/engine/internal/domain/entity/company"
	"github.com/haily-id/engine/internal/domain/repository"
	"gorm.io/gorm"
)

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) repository.CompanyRepository {
	return &companyRepository{
		db: db,
	}
}

func (r *companyRepository) Create(ctx context.Context, c *company.Company) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *companyRepository) FindByID(ctx context.Context, id int64) (*company.Company, error) {
	var c company.Company
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("company not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *companyRepository) FindByCode(ctx context.Context, code string) (*company.Company, error) {
	var c company.Company
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("company not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *companyRepository) Update(ctx context.Context, c *company.Company) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *companyRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&company.Company{}, id).Error
}

func (r *companyRepository) List(ctx context.Context, limit, offset int) ([]*company.Company, error) {
	var companies []*company.Company
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&companies).Error
	return companies, err
}
