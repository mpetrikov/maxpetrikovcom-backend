package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/config"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/postgres"
)

type App struct {
	config config.Config
	logger *slog.Logger

	db         *pgxpool.Pool
	httpServer *http.Server
}

func New(
	ctx context.Context,
	cfg config.Config,
	logger *slog.Logger,
) (*App, error) {
	db, err := postgres.NewPostgresPool(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"connect to PostgreSQL: %w",
			err,
		)
	}

	logger.Info("connected to PostgreSQL")

	deps := buildDependencies(
		cfg,
		db,
	)

	server := buildHTTPServer(
		cfg,
		logger,
		db,
		deps,
	)

	return &App{
		config:     cfg,
		logger:     logger,
		db:         db,
		httpServer: server,
	}, nil
}
