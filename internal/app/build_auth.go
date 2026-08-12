package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/config"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/httpserver"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/postgres"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/auth"
)

func buildTokenService(
	cfg config.Config,
) *auth.TokenService {
	return auth.NewTokenService(
		cfg.JWTSecret,
		cfg.JWTIssuer,
		cfg.JWTAccessTTL,
	)
}

func buildAuthHandler(
	cfg config.Config,
	db *pgxpool.Pool,
	tokenService *auth.TokenService,
) *httpserver.AuthHandler {
	userRepository := postgres.NewUserRepository(db)

	userIdentityRepository :=
		postgres.NewUserIdentityRepository(db)

	refreshSessionRepository :=
		postgres.NewRefreshSessionRepository(db)

	authService := auth.NewService(
		userRepository,
		userIdentityRepository,
	)

	googleOAuth := auth.NewGoogleOAuth(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
	)

	refreshTokenService := auth.NewRefreshTokenService(
		refreshSessionRepository,
		userRepository,
		tokenService,
		cfg.JWTRefreshTTL,
	)

	return httpserver.NewAuthHandler(
		authService,
		googleOAuth,
		tokenService,
		refreshTokenService,
	)
}
