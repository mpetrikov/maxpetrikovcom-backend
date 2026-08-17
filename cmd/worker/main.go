package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/config"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/worker"
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

	w, err := worker.New(
		context.Background(),
		cfg,
		logger,
	)
	if err != nil {
		logger.Error(
			"failed to initialize worker",
			"error", err,
		)
		os.Exit(1)
	}

	if err := w.Run(); err != nil {
		logger.Error(
			"worker stopped with error",
			"error", err,
		)
		os.Exit(1)
	}
}
