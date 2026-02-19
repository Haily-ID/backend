package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hibiken/asynq"
	userEntity "github.com/haily-id/engine/internal/domain/entity/user"
	"github.com/haily-id/engine/internal/domain/repository"
	"github.com/haily-id/engine/internal/pkg/asynq/tasks"
	"github.com/haily-id/engine/internal/pkg/mailer"
	"github.com/haily-id/engine/internal/pkg/snowflake"
	"golang.org/x/crypto/bcrypt"
)

// ─── Request DTOs ───────────────────────────────────────────────

type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name"     validate:"required,min=2"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
	OTP   string `json:"otp"   validate:"required,len=6"`
}

type ResendOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ─── Use Case ───────────────────────────────────────────────────

type UseCase struct {
	userRepo    repository.UserRepository
	evRepo      repository.EmailVerificationRepository
	mailer      mailer.Mailer
	asynqClient interface {
		Enqueue(task *asynq.Task, opts ...asynq.Option) error
	}
	jwtSecret      string
	jwtExpiryHours int
}

type Config struct {
	JWTSecret      string
	JWTExpiryHours int
}

func NewUseCase(
	userRepo repository.UserRepository,
	evRepo repository.EmailVerificationRepository,
	m mailer.Mailer,
	asynqClient interface {
		Enqueue(task *asynq.Task, opts ...asynq.Option) error
	},
	cfg Config,
) *UseCase {
	return &UseCase{
		userRepo:       userRepo,
		evRepo:         evRepo,
		mailer:         m,
		asynqClient:    asynqClient,
		jwtSecret:      cfg.JWTSecret,
		jwtExpiryHours: cfg.JWTExpiryHours,
	}
}

func (uc *UseCase) Register(ctx context.Context, req RegisterRequest) (*userEntity.User, error) {
	existing, _ := uc.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	id, err := snowflake.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ID: %w", err)
	}

	pwd := string(hashedPassword)
	u := &userEntity.User{
		ID:       id,
		Email:    req.Email,
		Password: &pwd,
		Name:     req.Name,
		Status:   userEntity.StatusPendingVerification,
	}

	if err := uc.userRepo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := uc.createAndSendOTP(ctx, u, userEntity.VerificationTypeEmailVerification); err != nil {
		return nil, err
	}

	return u, nil
}

func (uc *UseCase) Login(ctx context.Context, req LoginRequest) (*userEntity.User, string, error) {
	u, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if u.Status == userEntity.StatusPendingVerification {
		return nil, "", errors.New("email not verified")
	}

	if u.Status == userEntity.StatusSuspended {
		return nil, "", errors.New("account suspended")
	}

	if u.Password == nil {
		return nil, "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(req.Password)); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	token, err := uc.generateJWT(u)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	now := time.Now()
	u.LastLoginAt = &now
	_ = uc.userRepo.Update(ctx, u)

	return u, token, nil
}

func (uc *UseCase) VerifyEmail(ctx context.Context, req VerifyEmailRequest) (*userEntity.User, error) {
	ev, err := uc.evRepo.FindByToken(ctx, req.Token)
	if err != nil {
		return nil, errors.New("invalid or expired verification token")
	}

	if ev.IsUsed {
		return nil, errors.New("verification token already used")
	}

	if time.Now().After(ev.ExpiresAt) {
		return nil, errors.New("verification token expired")
	}

	if ev.AttemptsUsed >= ev.MaxAttempts {
		return nil, errors.New("max verification attempts exceeded")
	}

	if ev.OTPCode != req.OTP {
		_ = uc.evRepo.IncrementAttempts(ctx, ev.ID)
		return nil, errors.New("invalid OTP code")
	}

	if err := uc.evRepo.MarkUsed(ctx, ev.ID); err != nil {
		return nil, fmt.Errorf("failed to mark verification used: %w", err)
	}

	u, err := uc.userRepo.FindByID(ctx, ev.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	now := time.Now()
	u.Status = userEntity.StatusActive
	u.EmailVerifiedAt = &now

	if err := uc.userRepo.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to activate user: %w", err)
	}

	return u, nil
}

func (uc *UseCase) ResendOTP(ctx context.Context, req ResendOTPRequest) error {
	u, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	if u.Status == userEntity.StatusActive {
		return errors.New("email already verified")
	}

	if u.Status == userEntity.StatusSuspended {
		return errors.New("account suspended")
	}

	if err := uc.evRepo.InvalidateByUserIDAndType(ctx, u.ID, userEntity.VerificationTypeEmailVerification); err != nil {
		return fmt.Errorf("failed to invalidate previous OTPs: %w", err)
	}

	return uc.createAndSendOTP(ctx, u, userEntity.VerificationTypeEmailVerification)
}

func (uc *UseCase) GetMe(ctx context.Context, userID int64) (*userEntity.User, error) {
	return uc.userRepo.FindByID(ctx, userID)
}

// ─── Helpers ────────────────────────────────────────────────────

func (uc *UseCase) createAndSendOTP(ctx context.Context, u *userEntity.User, verType string) error {
	otp, err := generateOTP()
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	token, err := generateToken()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	evID, err := snowflake.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate ID: %w", err)
	}

	ev := &userEntity.EmailVerification{
		ID:          evID,
		UserID:      u.ID,
		Type:        verType,
		Token:       token,
		OTPCode:     otp,
		Email:       u.Email,
		MaxAttempts: userEntity.MaxOTPAttempts,
		ExpiresAt:   time.Now().Add(userEntity.OTPExpiry),
	}

	if err := uc.evRepo.Create(ctx, ev); err != nil {
		return fmt.Errorf("failed to create verification: %w", err)
	}

	task, err := tasks.NewSendOTPEmailTask(u.Email, u.Name, otp, verType)
	if err != nil {
		return fmt.Errorf("failed to create email task: %w", err)
	}

	if err := uc.asynqClient.Enqueue(task, asynq.Queue("default")); err != nil {
		return fmt.Errorf("failed to enqueue email task: %w", err)
	}

	return nil
}

func (uc *UseCase) generateJWT(u *userEntity.User) (string, error) {
	expiry := time.Duration(uc.jwtExpiryHours) * time.Hour
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"email":   u.Email,
		"status":  u.Status,
		"exp":     time.Now().Add(expiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

func generateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
