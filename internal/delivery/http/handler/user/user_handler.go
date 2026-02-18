package user

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/haily-id/engine/internal/domain/dto/user"
	"github.com/haily-id/engine/internal/pkg/response"
	"github.com/haily-id/engine/internal/pkg/validator"
	userUseCase "github.com/haily-id/engine/internal/usecase/user"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUC    *userUseCase.UserUseCase
	jwtSecret string
}

func NewUserHandler(userUC *userUseCase.UserUseCase, jwtSecret string) *UserHandler {
	return &UserHandler{
		userUC:    userUC,
		jwtSecret: jwtSecret,
	}
}

func (h *UserHandler) Register(c echo.Context) error {
	var req userUseCase.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}

	if err := validator.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}

	u, err := h.userUC.Register(c.Request().Context(), req)
	if err != nil {
		if err.Error() == "email already registered" {
			return response.Error(c, http.StatusConflict, response.ErrEmailAlreadyExists)
		}
		return response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
	}

	return response.Created(c, user.ToDTO(u))
}

func (h *UserHandler) Login(c echo.Context) error {
	var req userUseCase.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}

	if err := validator.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}

	u, err := h.userUC.Login(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, response.ErrInvalidCredentials)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID,
		"email":   u.Email,
		"role":    u.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
	}

	return response.Success(c, map[string]interface{}{
		"token": tokenString,
		"user":  user.ToDTO(u),
	})
}

func (h *UserHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	u, err := h.userUC.GetByID(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusNotFound, response.ErrUserNotFound)
	}

	return response.Success(c, user.ToDTO(u))
}

func (h *UserHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrInvalidUserID)
	}

	u, err := h.userUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, response.ErrUserNotFound)
	}

	return response.Success(c, user.ToDTO(u))
}

func (h *UserHandler) GetWithCompanies(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	u, err := h.userUC.GetByIDWithCompanies(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusNotFound, response.ErrUserNotFound)
	}

	return response.Success(c, user.ToUserWithCompaniesDTO(u))
}

func (h *UserHandler) Update(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	var req userUseCase.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}

	if err := validator.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}

	u, err := h.userUC.Update(c.Request().Context(), userID, req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, response.ErrUserUpdateFailed)
	}

	return response.Success(c, user.ToDTO(u))
}

func (h *UserHandler) Delete(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	if err := h.userUC.Delete(c.Request().Context(), userID); err != nil {
		return response.Error(c, http.StatusInternalServerError, response.ErrUserDeleteFailed)
	}

	return response.NoContent(c)
}

func (h *UserHandler) List(c echo.Context) error {
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

	users, err := h.userUC.List(c.Request().Context(), limit, offset)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, response.ErrUserListFailed)
	}

	dtos := make([]interface{}, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, user.ToDTO(u))
	}

	return response.Success(c, dtos)
}

func (h *UserHandler) JoinCompany(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	companyIDStr := c.Param("company_id")
	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrInvalidCompanyID)
	}

	if err := h.userUC.JoinCompany(c.Request().Context(), userID, companyID); err != nil {
		if err.Error() == "company not found" {
			return response.Error(c, http.StatusNotFound, response.ErrCompanyNotFound)
		}
		return response.Error(c, http.StatusInternalServerError, response.ErrJoinCompanyFailed)
	}

	return response.NoContent(c)
}

func (h *UserHandler) LeaveCompany(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	companyIDStr := c.Param("company_id")
	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrInvalidCompanyID)
	}

	if err := h.userUC.LeaveCompany(c.Request().Context(), userID, companyID); err != nil {
		return response.Error(c, http.StatusInternalServerError, response.ErrLeaveCompanyFailed)
	}

	return response.NoContent(c)
}
