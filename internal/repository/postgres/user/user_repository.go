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
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).First(&u, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByIDWithCompanies(ctx context.Context, id int64) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Preload("Companies").First(&u, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&user.User{}, id).Error
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	var users []*user.User
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *userRepository) AddCompany(ctx context.Context, userID, companyID int64) error {
	userCompany := &user.UserCompany{
		UserID:    userID,
		CompanyID: companyID,
	}
	return r.db.WithContext(ctx).Create(userCompany).Error
}

func (r *userRepository) RemoveCompany(ctx context.Context, userID, companyID int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND company_id = ?", userID, companyID).
		Delete(&user.UserCompany{}).Error
}

func (r *userRepository) GetCompaniesByUserID(ctx context.Context, userID int64) ([]*user.Company, error) {
	var companies []*user.Company
	err := r.db.WithContext(ctx).
		Table("companies").
		Joins("JOIN user_companies ON user_companies.company_id = companies.id").
		Where("user_companies.user_id = ?", userID).
		Find(&companies).Error
	return companies, err
}
