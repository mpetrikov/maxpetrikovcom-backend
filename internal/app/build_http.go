package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/config"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/httpserver"
)

func buildHTTPServer(
	cfg config.Config,
	logger *slog.Logger,
	db *pgxpool.Pool,
	deps dependencies,
) *http.Server {
	api := httpserver.New(
		logger,
		db,
		deps.authHandler,
		deps.userHandler,
		deps.tokenService,
		deps.labHandler,
	)

	return &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           api.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
