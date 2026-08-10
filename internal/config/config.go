package config

import (
	"fmt"
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

	JWTSecret    string
	JWTIssuer    string
	JWTAccessTTL time.Duration

	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	jwtAccessTTL, err := time.ParseDuration(
		getEnv("JWT_ACCESS_TTL", "15m"),
	)
	if err != nil {
		return Config{}, fmt.Errorf(
			"parse JWT_ACCESS_TTL: %w",
			err,
		)
	}

	return Config{
		Environment:        getEnv("APP_ENV", "local"),
		HTTPAddress:        getEnv("HTTP_ADDRESS", ":8080"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		RabbitMQURL:        getEnv("RABBITMQ_URL", ""),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),

		GoogleRedirectURL: getEnv("GOOGLE_REDIRECT_URL", ""),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		JWTIssuer:         getEnv("JWT_ISSUER", "maxpetrikov.com"),
		JWTAccessTTL:      jwtAccessTTL,

		ShutdownTimeout: 10 * time.Second,
	}, nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
