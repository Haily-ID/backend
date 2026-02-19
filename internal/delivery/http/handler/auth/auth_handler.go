package auth

import (
	"net/http"

	authDTO "github.com/haily-id/engine/internal/domain/dto/auth"
	userDTO "github.com/haily-id/engine/internal/domain/dto/user"
	"github.com/haily-id/engine/internal/pkg/response"
	"github.com/haily-id/engine/internal/pkg/validator"
	"github.com/haily-id/engine/internal/usecase/auth"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	authUC *auth.UseCase
}

func NewHandler(authUC *auth.UseCase) *Handler {
	return &Handler{authUC: authUC}
}

func (h *Handler) Register(c echo.Context) error {
	var req auth.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}
	if err := validator.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}

	u, err := h.authUC.Register(c.Request().Context(), req)
	if err != nil {
		if err.Error() == "email already registered" {
			return response.Error(c, http.StatusConflict, response.ErrEmailAlreadyExists)
		}
		return response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
	}

	return response.Created(c, authDTO.RegisterResponse{
		User:    userDTO.ToDTO(u),
		Message: "Registration successful. Please check your email for OTP verification.",
	})
}

func (h *Handler) Login(c echo.Context) error {
	var req auth.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}
	if err := validator.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}

	u, token, err := h.authUC.Login(c.Request().Context(), req)
	if err != nil {
		switch err.Error() {
		case "invalid email or password":
			return response.Error(c, http.StatusUnauthorized, response.ErrInvalidCredentials)
		case "email not verified":
			return response.Error(c, http.StatusForbidden, response.ErrEmailNotVerified)
		case "account suspended":
			return response.Error(c, http.StatusForbidden, response.ErrAccountSuspended)
		}
		return response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
	}

	return response.Success(c, authDTO.LoginResponse{
		Token: token,
		User:  userDTO.ToDTO(u),
	})
}

func (h *Handler) VerifyEmail(c echo.Context) error {
	var req auth.VerifyEmailRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}
	if err := validator.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}

	u, err := h.authUC.VerifyEmail(c.Request().Context(), req)
	if err != nil {
		switch err.Error() {
		case "invalid or expired verification token":
			return response.Error(c, http.StatusBadRequest, response.ErrInvalidVerificationToken)
		case "verification token already used":
			return response.Error(c, http.StatusBadRequest, response.ErrVerificationTokenUsed)
		case "verification token expired":
			return response.Error(c, http.StatusBadRequest, response.ErrVerificationTokenExpired)
		case "max verification attempts exceeded":
			return response.Error(c, http.StatusTooManyRequests, response.ErrMaxOTPAttemptsExceeded)
		case "invalid OTP code":
			return response.Error(c, http.StatusBadRequest, response.ErrInvalidOTP)
		}
		return response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
	}

	return response.Success(c, userDTO.ToDTO(u))
}

func (h *Handler) ResendOTP(c echo.Context) error {
	var req auth.ResendOTPRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}
	if err := validator.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, response.ErrValidation)
	}

	err := h.authUC.ResendOTP(c.Request().Context(), req)
	if err != nil {
		switch err.Error() {
		case "user not found":
			return response.Error(c, http.StatusNotFound, response.ErrUserNotFound)
		case "email already verified":
			return response.Error(c, http.StatusConflict, response.ErrEmailAlreadyVerified)
		case "account suspended":
			return response.Error(c, http.StatusForbidden, response.ErrAccountSuspended)
		}
		return response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
	}

	return response.Success(c, map[string]string{
		"message": "OTP sent successfully. Please check your email.",
	})
}

func (h *Handler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	u, err := h.authUC.GetMe(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusNotFound, response.ErrUserNotFound)
	}

	return response.Success(c, userDTO.ToDTO(u))
}
