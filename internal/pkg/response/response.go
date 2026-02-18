package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	ErrValidation         = "VALIDATION_ERROR"
	ErrUnauthorized       = "UNAUTHORIZED"
	ErrForbidden          = "FORBIDDEN"
	ErrInternalServer     = "INTERNAL_SERVER_ERROR"
	ErrServiceUnavailable = "SERVICE_UNAVAILABLE"

	ErrUserNotFound       = "USER_NOT_FOUND"
	ErrInvalidCredentials = "INVALID_CREDENTIALS"
	ErrEmailAlreadyExists = "EMAIL_ALREADY_EXISTS"
	ErrInvalidEmail       = "INVALID_EMAIL"
	ErrInvalidPassword    = "INVALID_PASSWORD"
	ErrInvalidUserID      = "INVALID_USER_ID"
	ErrUserUpdateFailed   = "USER_UPDATE_FAILED"
	ErrUserDeleteFailed   = "USER_DELETE_FAILED"
	ErrUserListFailed     = "USER_LIST_FAILED"

	ErrCompanyNotFound          = "COMPANY_NOT_FOUND"
	ErrCompanyCodeAlreadyExists = "COMPANY_CODE_ALREADY_EXISTS"
	ErrInvalidCompanyID         = "INVALID_COMPANY_ID"
	ErrInvalidCompanyCode       = "INVALID_COMPANY_CODE"
	ErrCompanyCreateFailed      = "COMPANY_CREATE_FAILED"
	ErrCompanyUpdateFailed      = "COMPANY_UPDATE_FAILED"
	ErrCompanyDeleteFailed      = "COMPANY_DELETE_FAILED"
	ErrCompanyListFailed        = "COMPANY_LIST_FAILED"
	ErrJoinCompanyFailed        = "JOIN_COMPANY_FAILED"
	ErrLeaveCompanyFailed       = "LEAVE_COMPANY_FAILED"
	ErrAlreadyCompanyMember     = "ALREADY_COMPANY_MEMBER"
	ErrNotCompanyMember         = "NOT_COMPANY_MEMBER"
	ErrCannotLeaveOwnCompany    = "CANNOT_LEAVE_OWN_COMPANY"
)

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, SuccessResponse{
		Data: data,
	})
}

func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, SuccessResponse{
		Data: data,
	})
}

func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func Error(c echo.Context, statusCode int, errorCode string) error {
	return c.JSON(statusCode, ErrorResponse{
		Error: errorCode,
	})
}
