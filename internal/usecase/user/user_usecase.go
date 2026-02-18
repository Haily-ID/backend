package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/haily-id/engine/internal/domain/entity/user"
	"github.com/haily-id/engine/internal/domain/repository"
	"github.com/haily-id/engine/internal/pkg/snowflake"
	"github.com/haily-id/engine/internal/repository/redis"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo    repository.UserRepository
	companyRepo repository.CompanyRepository
	cache       *redis.Cache
}

func NewUserUseCase(
	userRepo repository.UserRepository,
	companyRepo repository.CompanyRepository,
	cache *redis.Cache,
) *UserUseCase {
	return &UserUseCase{
		userRepo:    userRepo,
		companyRepo: companyRepo,
		cache:       cache,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

func (uc *UserUseCase) Register(ctx context.Context, req RegisterRequest) (*user.User, error) {
	// Check if user already exists
	existingUser, _ := uc.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate Snowflake ID
	id, err := snowflake.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ID: %w", err)
	}

	// Create user
	newUser := &user.User{
		ID:       id,
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     "user",
	}

	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Cache the user
	_ = uc.cache.Set(ctx, redis.UserKey(newUser.ID), newUser, time.Hour)

	return newUser, nil
}

func (uc *UserUseCase) Login(ctx context.Context, req LoginRequest) (*user.User, error) {
	// Find user by email
	u, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return u, nil
}

func (uc *UserUseCase) GetByID(ctx context.Context, id int64) (*user.User, error) {
	// Try cache first
	var u user.User
	err := uc.cache.Get(ctx, redis.UserKey(id), &u)
	if err == nil {
		return &u, nil
	}

	// Cache miss, get from database
	u2, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = uc.cache.Set(ctx, redis.UserKey(id), u2, time.Hour)

	return u2, nil
}

func (uc *UserUseCase) GetByIDWithCompanies(ctx context.Context, id int64) (*user.User, error) {
	return uc.userRepo.FindByIDWithCompanies(ctx, id)
}

func (uc *UserUseCase) Update(ctx context.Context, id int64, req UpdateUserRequest) (*user.User, error) {
	// Get existing user
	u, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	u.Name = req.Name

	// Save to database
	if err := uc.userRepo.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Invalidate cache
	_ = uc.cache.Delete(ctx, redis.UserKey(id))

	return u, nil
}

func (uc *UserUseCase) Delete(ctx context.Context, id int64) error {
	if err := uc.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Invalidate cache
	_ = uc.cache.Delete(ctx, redis.UserKey(id))

	return nil
}

func (uc *UserUseCase) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	return uc.userRepo.List(ctx, limit, offset)
}

func (uc *UserUseCase) JoinCompany(ctx context.Context, userID, companyID int64) error {
	// Verify user exists
	_, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify company exists
	_, err = uc.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return errors.New("company not found")
	}

	// Add user to company
	if err := uc.userRepo.AddCompany(ctx, userID, companyID); err != nil {
		return fmt.Errorf("failed to join company: %w", err)
	}

	return nil
}

func (uc *UserUseCase) LeaveCompany(ctx context.Context, userID, companyID int64) error {
	if err := uc.userRepo.RemoveCompany(ctx, userID, companyID); err != nil {
		return fmt.Errorf("failed to leave company: %w", err)
	}

	return nil
}

func (uc *UserUseCase) GetUserCompanies(ctx context.Context, userID int64) ([]*user.Company, error) {
	return uc.userRepo.GetCompaniesByUserID(ctx, userID)
}
