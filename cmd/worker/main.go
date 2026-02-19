package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	pkgAsynq "github.com/haily-id/engine/internal/pkg/asynq"
	"github.com/haily-id/engine/internal/pkg/asynq/tasks"
	"github.com/haily-id/engine/internal/pkg/config"
	"github.com/haily-id/engine/internal/pkg/database"
	"github.com/haily-id/engine/internal/pkg/logger"
	"github.com/haily-id/engine/internal/pkg/mailer"
	"github.com/haily-id/engine/internal/pkg/snowflake"
	asynqLib "github.com/hibiken/asynq"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	logger.Init("WORKER")

	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := snowflake.Init(cfg.Snowflake.MachineID); err != nil {
		log.Fatalf("Failed to initialize Snowflake: %v", err)
	}

	_, err = database.NewPostgresDB(database.Config{
		DSN:             cfg.Database.DSN(),
		MaxOpenConns:    10,
		MaxIdleConns:    2,
		ConnMaxLifetime: 30 * time.Minute,
		LogLevel:        gormLogger.Warn,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	m := mailer.New(mailer.Config{
		Driver:   cfg.Mailer.Driver,
		From:     cfg.Mailer.From,
		Host:     cfg.Mailer.Host,
		Port:     cfg.Mailer.Port,
		Username: cfg.Mailer.Username,
		Password: cfg.Mailer.Password,
	})

	server := pkgAsynq.NewServer(cfg.Asynq.RedisAddr, 10)

	mux := asynqLib.NewServeMux()
	mux.HandleFunc(tasks.TypeSendOTPEmail, handleSendOTPEmail(m))

	logger.Info("Starting worker...")

	go func() {
		if err := server.Start(mux); err != nil {
			logger.Errorf("Worker error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker...")
	server.Shutdown()
	logger.Info("Worker stopped")
}

func handleSendOTPEmail(m mailer.Mailer) asynqLib.HandlerFunc {
	return func(ctx context.Context, t *asynqLib.Task) error {
		var payload tasks.SendOTPEmailPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}
		logger.Infof("Sending OTP email to %s", payload.To)
		return m.SendOTP(payload.To, payload.Name, payload.OTP, payload.Purpose)
	}
}
