package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	authHandler "github.com/haily-id/engine/internal/delivery/http/handler/auth"
	"github.com/haily-id/engine/internal/delivery/http/route"
	userEntity "github.com/haily-id/engine/internal/domain/entity/user"
	pkgAsynq "github.com/haily-id/engine/internal/pkg/asynq"
	"github.com/haily-id/engine/internal/pkg/config"
	"github.com/haily-id/engine/internal/pkg/database"
	"github.com/haily-id/engine/internal/pkg/logger"
	"github.com/haily-id/engine/internal/pkg/mailer"
	"github.com/haily-id/engine/internal/pkg/snowflake"
	"github.com/haily-id/engine/internal/pkg/validator"
	userRepo "github.com/haily-id/engine/internal/repository/postgres/user"
	authUC "github.com/haily-id/engine/internal/usecase/auth"
	"github.com/labstack/echo/v4"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	logger.Init("API")

	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := snowflake.Init(cfg.Snowflake.MachineID); err != nil {
		log.Fatalf("Failed to initialize Snowflake: %v", err)
	}

	validator.Init()

	db, err := database.NewPostgresDB(database.Config{
		DSN:             cfg.Database.DSN(),
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30 * time.Minute,
		LogLevel:        gormLogger.Info,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	if err := db.AutoMigrate(
		&userEntity.User{},
		&userEntity.EmailVerification{},
	); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	logger.Info("Database migration completed")

	asynqClient := pkgAsynq.NewClient(cfg.Asynq.RedisAddr)
	defer asynqClient.Close()

	m := mailer.New(mailer.Config{
		Driver:   cfg.Mailer.Driver,
		FromName: cfg.Mailer.FromName,
		From:     cfg.Mailer.From,
		Host:     cfg.Mailer.Host,
		Port:     cfg.Mailer.Port,
		Username: cfg.Mailer.Username,
		Password: cfg.Mailer.Password,
	})

	userRepository := userRepo.NewUserRepository(db)
	evRepository := userRepo.NewEmailVerificationRepository(db)

	authUseCase := authUC.NewUseCase(
		userRepository,
		evRepository,
		m,
		asynqClient,
		authUC.Config{
			JWTSecret:      cfg.JWT.Secret,
			JWTExpiryHours: cfg.JWT.ExpirationHour,
		},
	)

	authH := authHandler.NewHandler(authUseCase)

	e := echo.New()
	e.HideBanner = true

	route.Setup(e, route.RouteConfig{
		AuthHandler: authH,
		JWTSecret:   cfg.JWT.Secret,
	})

	go func() {
		addr := fmt.Sprintf(":%s", cfg.App.Port)
		logger.Infof("Starting API server on %s", addr)
		if err := e.Start(addr); err != nil {
			logger.Errorf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Errorf("Server shutdown error: %v", err)
	}

	logger.Info("Server stopped")
}
