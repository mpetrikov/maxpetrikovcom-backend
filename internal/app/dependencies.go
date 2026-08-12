package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/config"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/httpserver"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/auth"
)

type dependencies struct {
	authHandler  *httpserver.AuthHandler
	userHandler  *httpserver.UserHandler
	labHandler   *httpserver.LabHandler
	tokenService *auth.TokenService
}

func buildDependencies(
	cfg config.Config,
	db *pgxpool.Pool,
) dependencies {
	tokenService := buildTokenService(cfg)

	return dependencies{
		authHandler: buildAuthHandler(
			cfg,
			db,
			tokenService,
		),
		userHandler:  buildUserHandler(db),
		labHandler:   buildLabHandler(db),
		tokenService: tokenService,
	}
}
