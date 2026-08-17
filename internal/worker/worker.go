package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/config"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/queue/rabbitmq"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/postgres"
	labsessionservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/labsession"
)

type Worker struct {
	config config.Config
	logger *slog.Logger

	db       *pgxpool.Pool
	rabbitMQ *rabbitmq.Client

	labSessionService *labsessionservice.Service
}

func New(
	ctx context.Context,
	cfg config.Config,
	logger *slog.Logger,
) (*Worker, error) {
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

	rabbitMQ, err := rabbitmq.New(
		cfg.RabbitMQURL,
	)
	if err != nil {
		db.Close()

		return nil, fmt.Errorf(
			"connect to RabbitMQ: %w",
			err,
		)
	}

	labRepository := postgres.NewLabRepository(db)
	labSessionRepository := postgres.NewLabSessionRepository(db)

	labSessionService := labsessionservice.NewService(
		labRepository,
		labSessionRepository,
		rabbitMQ,
	)

	return &Worker{
		config:            cfg,
		logger:            logger,
		db:                db,
		rabbitMQ:          rabbitMQ,
		labSessionService: labSessionService,
	}, nil
}
