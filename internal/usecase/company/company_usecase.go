package company

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/haily-id/engine/internal/domain/entity/company"
	"github.com/haily-id/engine/internal/domain/repository"
	"github.com/haily-id/engine/internal/pkg/snowflake"
	"github.com/haily-id/engine/internal/repository/redis"
)

type CompanyUseCase struct {
	companyRepo repository.CompanyRepository
	cache       *redis.Cache
}

func NewCompanyUseCase(
	companyRepo repository.CompanyRepository,
	cache *redis.Cache,
) *CompanyUseCase {
	return &CompanyUseCase{
		companyRepo: companyRepo,
		cache:       cache,
	}
}

type CreateCompanyRequest struct {
	Name    string `json:"name" validate:"required,min=2"`
	Code    string `json:"code" validate:"required,min=2,max=50"`
	Address string `json:"address"`
}

type UpdateCompanyRequest struct {
	Name    string `json:"name" validate:"required,min=2"`
	Address string `json:"address"`
}

func (uc *CompanyUseCase) Create(ctx context.Context, req CreateCompanyRequest) (*company.Company, error) {
	// Check if code already exists
	existingCompany, _ := uc.companyRepo.FindByCode(ctx, req.Code)
	if existingCompany != nil {
		return nil, errors.New("company code already exists")
	}

	// Generate Snowflake ID
	id, err := snowflake.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ID: %w", err)
	}

	// Create company
	newCompany := &company.Company{
		ID:      id,
		Name:    req.Name,
		Code:    req.Code,
		Address: req.Address,
	}

	if err := uc.companyRepo.Create(ctx, newCompany); err != nil {
		return nil, fmt.Errorf("failed to create company: %w", err)
	}

	// Cache the company
	_ = uc.cache.Set(ctx, redis.CompanyKey(newCompany.ID), newCompany, 24*time.Hour)
	_ = uc.cache.Set(ctx, redis.CompanyCodeKey(newCompany.Code), newCompany, 24*time.Hour)

	return newCompany, nil
}

func (uc *CompanyUseCase) GetByID(ctx context.Context, id int64) (*company.Company, error) {
	// Try cache first
	var c company.Company
	err := uc.cache.Get(ctx, redis.CompanyKey(id), &c)
	if err == nil {
		return &c, nil
	}

	// Cache miss, get from database
	c2, err := uc.companyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = uc.cache.Set(ctx, redis.CompanyKey(id), c2, 24*time.Hour)

	return c2, nil
}

func (uc *CompanyUseCase) GetByCode(ctx context.Context, code string) (*company.Company, error) {
	// Try cache first
	var c company.Company
	err := uc.cache.Get(ctx, redis.CompanyCodeKey(code), &c)
	if err == nil {
		return &c, nil
	}

	// Cache miss, get from database
	c2, err := uc.companyRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = uc.cache.Set(ctx, redis.CompanyCodeKey(code), c2, 24*time.Hour)

	return c2, nil
}

func (uc *CompanyUseCase) Update(ctx context.Context, id int64, req UpdateCompanyRequest) (*company.Company, error) {
	// Get existing company
	c, err := uc.companyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	c.Name = req.Name
	c.Address = req.Address

	// Save to database
	if err := uc.companyRepo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("failed to update company: %w", err)
	}

	// Invalidate cache
	_ = uc.cache.Delete(ctx, redis.CompanyKey(id))
	_ = uc.cache.Delete(ctx, redis.CompanyCodeKey(c.Code))

	return c, nil
}

func (uc *CompanyUseCase) Delete(ctx context.Context, id int64) error {
	// Get company to retrieve code for cache invalidation
	c, err := uc.companyRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.companyRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	// Invalidate cache
	_ = uc.cache.Delete(ctx, redis.CompanyKey(id))
	_ = uc.cache.Delete(ctx, redis.CompanyCodeKey(c.Code))

	return nil
}

func (uc *CompanyUseCase) List(ctx context.Context, limit, offset int) ([]*company.Company, error) {
	return uc.companyRepo.List(ctx, limit, offset)
}
