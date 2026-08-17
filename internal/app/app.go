package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/config"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/queue/rabbitmq"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/postgres"
)

type App struct {
	config config.Config
	logger *slog.Logger

	db       *pgxpool.Pool
	rabbitMQ *rabbitmq.Client

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

	rabbitMQ, err := buildRabbitMQ(
		cfg.RabbitMQURL,
	)
	if err != nil {
		db.Close()

		return nil, fmt.Errorf(
			"initialize RabbitMQ: %w",
			err,
		)
	}

	logger.Info("connected to RabbitMQ")

	deps := buildDependencies(
		cfg,
		db,
		rabbitMQ,
	)

	server := buildHTTPServer(
		cfg,
		logger,
		db,
		deps,
	)

	return &App{
		config: cfg,
		logger: logger,

		db:       db,
		rabbitMQ: rabbitMQ,

		httpServer: server,
	}, nil
}
