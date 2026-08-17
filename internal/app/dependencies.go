package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/config"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/httpserver"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/queue/rabbitmq"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/auth"
)

type dependencies struct {
	authHandler       *httpserver.AuthHandler
	userHandler       *httpserver.UserHandler
	labHandler        *httpserver.LabHandler
	labSessionHandler *httpserver.LabSessionHandler
	tokenService      *auth.TokenService
}

func buildDependencies(
	cfg config.Config,
	db *pgxpool.Pool,
	rabbitMQ *rabbitmq.Client,
) dependencies {
	tokenService := buildTokenService(cfg)

	return dependencies{
		authHandler: buildAuthHandler(
			cfg,
			db,
			tokenService,
		),
		userHandler:       buildUserHandler(db),
		labHandler:        buildLabHandler(db),
		labSessionHandler: buildLabSessionHandler(db, rabbitMQ),
		tokenService:      tokenService,
	}
}
