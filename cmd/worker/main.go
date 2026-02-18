package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/haily-id/engine/internal/pkg/asynq"
	"github.com/haily-id/engine/internal/pkg/config"
	"github.com/haily-id/engine/internal/pkg/database"
	"github.com/haily-id/engine/internal/pkg/logger"
	"github.com/haily-id/engine/internal/pkg/snowflake"
	asynqLib "github.com/hibiken/asynq"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	// Initialize logger
	logger.Init("WORKER")

	// Load configuration
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Snowflake ID generator
	if err := snowflake.Init(cfg.Snowflake.MachineID); err != nil {
		log.Fatalf("Failed to initialize Snowflake: %v", err)
	}

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

	// Initialize Asynq server
	server := asynq.NewServer(cfg.Asynq.RedisAddr, 10)

	// Create task mux
	mux := asynqLib.NewServeMux()

	// Register task handlers here
	// Example:
	// mux.HandleFunc("email:send", handleEmailTask)
	// mux.HandleFunc("report:generate", handleReportTask)

	logger.Info("Starting worker...")

	// Start worker in goroutine
	go func() {
		if err := server.Start(mux); err != nil {
			logger.Errorf("Worker error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker...")

	// Graceful shutdown
	server.Shutdown()

	logger.Info("Worker stopped")
}

// Example task handler
// func handleEmailTask(ctx context.Context, t *asynq.Task) error {
// 	var payload EmailPayload
// 	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
// 		return err
// 	}
//
// 	logger.Infof("Sending email to %s", payload.To)
// 	// Send email logic here
//
// 	return nil
// }
