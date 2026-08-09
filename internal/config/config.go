package config

import (
	"os"
	"time"
)

type Config struct {
	Environment        string
	HTTPAddress        string
	DatabaseURL        string
	RabbitMQURL        string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	ShutdownTimeout    time.Duration
}

func Load() Config {
	return Config{
		Environment:        getEnv("APP_ENV", "local"),
		HTTPAddress:        getEnv("HTTP_ADDRESS", ":8080"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		RabbitMQURL:        getEnv("RABBITMQ_URL", ""),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
		ShutdownTimeout:    10 * time.Second,
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
