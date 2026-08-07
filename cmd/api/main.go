package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/config"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/database"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/httpserver"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	if err := godotenv.Load(); err != nil {
		logger.Info(".env file not found, using environment variables")
	}

	cfg := config.Load()

	db, err := database.NewPostgresPool(
		context.Background(),
		cfg.DatabaseURL,
	)
	if err != nil {
		logger.Error(
			"failed to connect to PostgreSQL",
			"error", err,
		)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("connected to PostgreSQL")

	api := httpserver.New(logger, db)

	server := &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           api.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverError := make(chan error, 1)

	go func() {
		logger.Info(
			"starting HTTP server",
			"address", cfg.HTTPAddress,
			"environment", cfg.Environment,
		)

		serverError <- server.ListenAndServe()
	}()

	shutdownContext, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-serverError:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error(
				"HTTP server failed",
				"error", err,
			)
			os.Exit(1)
		}

	case <-shutdownContext.Done():
		logger.Info("shutdown signal received")

		timeoutContext, cancel := context.WithTimeout(
			context.Background(),
			cfg.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(timeoutContext); err != nil {
			logger.Error(
				"failed to shut down HTTP server",
				"error", err,
			)
			os.Exit(1)
		}

		logger.Info("HTTP server stopped")
	}
}
