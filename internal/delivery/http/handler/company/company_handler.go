package company

import (
	"net/http"
	"strconv"

	"github.com/haily-id/engine/internal/domain/dto/company"
	"github.com/haily-id/engine/internal/pkg/response"
	"github.com/haily-id/engine/internal/pkg/validator"
	companyUseCase "github.com/haily-id/engine/internal/usecase/company"
	"github.com/labstack/echo/v4"
)

type CompanyHandler struct {
	companyUC *companyUseCase.CompanyUseCase
}

func NewCompanyHandler(companyUC *companyUseCase.CompanyUseCase) *CompanyHandler {
	return &CompanyHandler{
		companyUC: companyUC,
	}
}

func (h *CompanyHandler) Create(c echo.Context) error {
	var req companyUseCase.CreateCompanyRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}

	if err := validator.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}

	comp, err := h.companyUC.Create(c.Request().Context(), req)
	if err != nil {
		if err.Error() == "company code already exists" {
			return response.Error(c, http.StatusConflict, response.ErrCompanyCodeAlreadyExists)
		}
		return response.Error(c, http.StatusInternalServerError, response.ErrCompanyCreateFailed)
	}

	return response.Created(c, company.ToDTO(comp))
}

func (h *CompanyHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrInvalidCompanyID)
	}

	comp, err := h.companyUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, response.ErrCompanyNotFound)
	}

	return response.Success(c, company.ToDTO(comp))
}

func (h *CompanyHandler) GetByCode(c echo.Context) error {
	code := c.Param("code")

	comp, err := h.companyUC.GetByCode(c.Request().Context(), code)
	if err != nil {
		return response.Error(c, http.StatusNotFound, response.ErrCompanyNotFound)
	}

	return response.Success(c, company.ToDTO(comp))
}

func (h *CompanyHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrInvalidCompanyID)
	}

	var req companyUseCase.UpdateCompanyRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}

	if err := validator.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}

	comp, err := h.companyUC.Update(c.Request().Context(), id, req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, response.ErrCompanyUpdateFailed)
	}

	return response.Success(c, company.ToDTO(comp))
}

func (h *CompanyHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrInvalidCompanyID)
	}

	if err := h.companyUC.Delete(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, response.ErrCompanyDeleteFailed)
	}

	return response.NoContent(c)
}

func (h *CompanyHandler) List(c echo.Context) error {
	limit := 20
	offset := 0

	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	companies, err := h.companyUC.List(c.Request().Context(), limit, offset)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, response.ErrCompanyListFailed)
	}

	dtos := make([]interface{}, 0, len(companies))
	for _, comp := range companies {
		dtos = append(dtos, company.ToDTO(comp))
	}

	return response.Success(c, dtos)
}
