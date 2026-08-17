package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/httpserver"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/postgres"

	queuecontracts "github.com/maxpetrikov/maxpetrikovcom-backend/internal/queue/contracts"
	labsessionservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/labsession"
)

func buildLabSessionHandler(
	db *pgxpool.Pool,
	publisher queuecontracts.LabJobPublisher,
) *httpserver.LabSessionHandler {
	labRepository :=
		postgres.NewLabRepository(db)

	labSessionRepository :=
		postgres.NewLabSessionRepository(db)

	sessionService :=
		labsessionservice.NewService(
			labRepository,
			labSessionRepository,
			publisher,
		)

	return httpserver.NewLabSessionHandler(
		sessionService,
	)
}
