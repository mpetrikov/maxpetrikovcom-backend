package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/httpserver"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/postgres"

	labservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/lab"
)

func buildLabHandler(
	db *pgxpool.Pool,
) *httpserver.LabHandler {
	labRepository :=
		postgres.NewLabRepository(db)

	labService :=
		labservice.NewService(labRepository)

	return httpserver.NewLabHandler(
		labService,
	)
}
