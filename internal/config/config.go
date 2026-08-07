package config

import (
	"os"
	"time"
)

type Config struct {
	Environment     string
	HTTPAddress     string
	DatabaseURL     string
	RabbitMQURL     string
	ShutdownTimeout time.Duration
}

func Load() Config {
	return Config{
		Environment:     getEnv("APP_ENV", "local"),
		HTTPAddress:     getEnv("HTTP_ADDRESS", ":8080"),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		RabbitMQURL:     getEnv("RABBITMQ_URL", ""),
		ShutdownTimeout: 10 * time.Second,
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
