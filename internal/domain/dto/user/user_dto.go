package user

import (
	"strconv"

	userEntity "github.com/haily-id/engine/internal/domain/entity/user"
)

type UserDTO struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type UserWithCompaniesDTO struct {
	ID        string       `json:"id"`
	Email     string       `json:"email"`
	Name      string       `json:"name"`
	Role      string       `json:"role"`
	CreatedAt int64        `json:"created_at"`
	UpdatedAt int64        `json:"updated_at"`
	Companies []CompanyDTO `json:"companies,omitempty"`
}

type CompanyDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Address   string `json:"address"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func ToDTO(user *userEntity.User) UserDTO {
	return UserDTO{
		ID:        strconv.FormatInt(user.ID, 10),
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}
}

func ToCompanyDTO(company *userEntity.Company) CompanyDTO {
	return CompanyDTO{
		ID:        strconv.FormatInt(company.ID, 10),
		Name:      company.Name,
		Code:      company.Code,
		Address:   company.Address,
		CreatedAt: company.CreatedAt.Unix(),
		UpdatedAt: company.UpdatedAt.Unix(),
	}
}

func ToUserWithCompaniesDTO(user *userEntity.User) UserWithCompaniesDTO {
	dto := UserWithCompaniesDTO{
		ID:        strconv.FormatInt(user.ID, 10),
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}

	if len(user.Companies) > 0 {
		dto.Companies = make([]CompanyDTO, len(user.Companies))
		for i, company := range user.Companies {
			dto.Companies[i] = ToCompanyDTO(&company)
		}
	}

	return dto
}
