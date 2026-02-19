package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App       AppConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	JWT       JWTConfig
	Snowflake SnowflakeConfig
	Asynq     AsynqConfig
	Mailer    MailerConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret         string
	ExpirationHour int
}

type SnowflakeConfig struct {
	MachineID int64
}

type AsynqConfig struct {
	RedisAddr string
}

type MailerConfig struct {
	Driver   string
	FromName string
	From     string
	Host     string
	Port     string
	Username string
	Password string
}

func Load(envFile string) (*Config, error) {
	_ = godotenv.Load(envFile)

	cfg := &Config{}

	cfg.App.Name = getEnv("APP_NAME", "Haily Backend")
	cfg.App.Env = getEnv("APP_ENV", "development")
	cfg.App.Port = getEnv("APP_PORT", "8080")

	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.Database.DBName = getEnv("DB_NAME", "haily")
	cfg.Database.SSLMode = getEnv("DB_SSL_MODE", "disable")

	cfg.Redis.Host = getEnv("REDIS_HOST", "localhost")
	cfg.Redis.Port = getEnv("REDIS_PORT", "6379")
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	if db, err := strconv.Atoi(getEnv("REDIS_DB", "0")); err == nil {
		cfg.Redis.DB = db
	}

	cfg.JWT.Secret = getEnv("JWT_SECRET", "change-me-in-production")
	if hours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOUR", "24")); err == nil {
		cfg.JWT.ExpirationHour = hours
	}

	if id, err := strconv.ParseInt(getEnv("SNOWFLAKE_MACHINE_ID", "1"), 10, 64); err == nil {
		cfg.Snowflake.MachineID = id
	}

	cfg.Asynq.RedisAddr = getEnv("ASYNQ_REDIS_ADDR", fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port))

	cfg.Mailer.Driver = getEnv("MAIL_DRIVER", "console")
	cfg.Mailer.FromName = getEnv("MAIL_FROM_NAME", "Haily")
	cfg.Mailer.From = getEnv("MAIL_FROM_ADDRESS", "noreply@haily.id")
	cfg.Mailer.Host = getEnv("MAIL_HOST", "")
	cfg.Mailer.Port = getEnv("MAIL_PORT", "587")
	cfg.Mailer.Username = getEnv("MAIL_USERNAME", "")
	cfg.Mailer.Password = getEnv("MAIL_PASSWORD", "")

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
