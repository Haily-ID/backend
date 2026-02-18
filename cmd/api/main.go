package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/haily-id/engine/internal/delivery/http/handler/company"
	"github.com/haily-id/engine/internal/delivery/http/handler/user"
	"github.com/haily-id/engine/internal/delivery/http/route"
	companyEntity "github.com/haily-id/engine/internal/domain/entity/company"
	userEntity "github.com/haily-id/engine/internal/domain/entity/user"
	"github.com/haily-id/engine/internal/pkg/config"
	"github.com/haily-id/engine/internal/pkg/database"
	"github.com/haily-id/engine/internal/pkg/logger"
	"github.com/haily-id/engine/internal/pkg/snowflake"
	"github.com/haily-id/engine/internal/pkg/validator"
	companyRepo "github.com/haily-id/engine/internal/repository/postgres/company"
	userRepo "github.com/haily-id/engine/internal/repository/postgres/user"
	"github.com/haily-id/engine/internal/repository/redis"
	companyUC "github.com/haily-id/engine/internal/usecase/company"
	userUC "github.com/haily-id/engine/internal/usecase/user"
	"github.com/labstack/echo/v4"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	// Initialize logger
	logger.Init("API")

	// Load configuration
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Snowflake ID generator
	if err := snowflake.Init(cfg.Snowflake.MachineID); err != nil {
		log.Fatalf("Failed to initialize Snowflake: %v", err)
	}

	// Initialize validator
	validator.Init()

	// Connect to PostgreSQL
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

	// Auto migrate
	if err := db.AutoMigrate(
		&userEntity.User{},
		&companyEntity.Company{},
		&userEntity.UserCompany{},
	); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	logger.Info("Database migration completed")

	// Connect to Redis
	cache, err := redis.NewCache(cfg.Redis.Addr(), cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer cache.Close()

	// Initialize repositories
	userRepository := userRepo.NewUserRepository(db)
	companyRepository := companyRepo.NewCompanyRepository(db)

	// Initialize use cases
	userUseCase := userUC.NewUserUseCase(userRepository, companyRepository, cache)
	companyUseCase := companyUC.NewCompanyUseCase(companyRepository, cache)

	// Initialize handlers
	userHandler := user.NewUserHandler(userUseCase, cfg.JWT.Secret)
	companyHandler := company.NewCompanyHandler(companyUseCase)

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true

	// Setup routes
	route.Setup(e, route.RouteConfig{
		UserHandler:    userHandler,
		CompanyHandler: companyHandler,
		JWTSecret:      cfg.JWT.Secret,
	})

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.App.Port)
		logger.Infof("Starting API server on %s", addr)
		if err := e.Start(addr); err != nil {
			logger.Errorf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Errorf("Server shutdown error: %v", err)
	}

	logger.Info("Server stopped")
}
