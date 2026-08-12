package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/httpserver"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/postgres"

	userservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/user"
)

func buildUserHandler(
	db *pgxpool.Pool,
) *httpserver.UserHandler {
	userRepository :=
		postgres.NewUserRepository(db)

	userService :=
		userservice.NewService(userRepository)

	return httpserver.NewUserHandler(
		userService,
	)
}
