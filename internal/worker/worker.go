package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/config"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/queue/rabbitmq"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/postgres"
	runnercontracts "github.com/maxpetrikov/maxpetrikovcom-backend/internal/runner/contracts"
	runnerfake "github.com/maxpetrikov/maxpetrikovcom-backend/internal/runner/fake"
	labservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/lab"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/labexecution"
	labsessionservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/labsession"
)

type Worker struct {
	config config.Config
	logger *slog.Logger

	db       *pgxpool.Pool
	rabbitMQ *rabbitmq.Client

	labSessionService   *labsessionservice.Service
	labService          *labservice.Service
	labRunner           runnercontracts.LabRunner
	labExecutionService *labexecution.Service
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

	labService := labservice.NewService(
		labRepository,
	)

	labSessionRepository := postgres.NewLabSessionRepository(db)

	labSessionService := labsessionservice.NewService(
		labRepository,
		labSessionRepository,
		rabbitMQ,
	)

	labRunner := runnerfake.NewLabRunner(logger)

	labExecutionService := labexecution.NewService(
		labService,
		labSessionService,
		labRunner,
		logger,
	)

	return &Worker{
		config:              cfg,
		logger:              logger,
		db:                  db,
		rabbitMQ:            rabbitMQ,
		labService:          labService,
		labSessionService:   labSessionService,
		labRunner:           labRunner,
		labExecutionService: labExecutionService,
	}, nil
}
