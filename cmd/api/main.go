package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/app"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/config"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	if err := godotenv.Load(); err != nil {
		logger.Info(".env file not found, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Error(
			"failed to load configuration",
			"error", err,
		)
		os.Exit(1)
	}

	application, err := app.New(
		context.Background(),
		cfg,
		logger,
	)
	if err != nil {
		logger.Error(
			"failed to initialize application",
			"error", err,
		)
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		logger.Error(
			"application stopped with error",
			"error", err,
		)
		os.Exit(1)
	}
}
