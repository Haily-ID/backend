package company

import (
	"strconv"

	companyEntity "github.com/haily-id/engine/internal/domain/entity/company"
)

type CompanyDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Address   string `json:"address"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func ToDTO(company *companyEntity.Company) CompanyDTO {
	return CompanyDTO{
		ID:        strconv.FormatInt(company.ID, 10),
		Name:      company.Name,
		Code:      company.Code,
		Address:   company.Address,
		CreatedAt: company.CreatedAt.Unix(),
		UpdatedAt: company.UpdatedAt.Unix(),
	}
}
