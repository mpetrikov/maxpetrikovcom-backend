package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/httpserver"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/postgres"

	labsessionservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/labsession"
)

func buildLabSessionHandler(
	db *pgxpool.Pool,
) *httpserver.LabSessionHandler {
	labRepository :=
		postgres.NewLabRepository(db)

	labSessionRepository :=
		postgres.NewLabSessionRepository(db)

	sessionService :=
		labsessionservice.NewService(
			labRepository,
			labSessionRepository,
		)

	return httpserver.NewLabSessionHandler(
		sessionService,
	)
}
