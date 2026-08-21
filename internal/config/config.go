package config

import (
	"fmt"
	"os"
	"time"
)

const (
	LabRunnerTypeFake       = "fake"
	LabRunnerTypeKubernetes = "kubernetes"
)

type Config struct {
	Environment        string
	HTTPAddress        string
	DatabaseURL        string
	RabbitMQURL        string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	JWTSecret     string
	JWTIssuer     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration

	LabRunnerType             string
	KubeconfigPath            string
	KubernetesLabNamespace    string
	KubernetesPodReadyTimeout time.Duration

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

	jwtRefreshTTL, err := time.ParseDuration(
		getEnv("JWT_REFRESH_TTL", "720h"),
	)
	if err != nil {
		return Config{}, fmt.Errorf(
			"parse JWT_REFRESH_TTL: %w",
			err,
		)
	}

	kubernetesPodReadyTimeout, err := time.ParseDuration(
		getEnv("KUBERNETES_POD_READY_TIMEOUT", "60s"),
	)
	if err != nil {
		return Config{}, fmt.Errorf(
			"parse KUBERNETES_POD_READY_TIMEOUT: %w",
			err,
		)
	}

	labRunnerType := getEnv("LAB_RUNNER_TYPE", LabRunnerTypeFake)
	if labRunnerType != LabRunnerTypeFake &&
		labRunnerType != LabRunnerTypeKubernetes {
		return Config{}, fmt.Errorf(
			"unsupported LAB_RUNNER_TYPE %q, expected %q or %q",
			labRunnerType,
			LabRunnerTypeFake,
			LabRunnerTypeKubernetes,
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
		JWTRefreshTTL:     jwtRefreshTTL,

		LabRunnerType:             labRunnerType,
		KubeconfigPath:            getEnv("KUBECONFIG_PATH", getEnv("KUBECONFIG", "")),
		KubernetesLabNamespace:    getEnv("KUBERNETES_LAB_NAMESPACE", "maxpetrikov-labs"),
		KubernetesPodReadyTimeout: kubernetesPodReadyTimeout,

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
