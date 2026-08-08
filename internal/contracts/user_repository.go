package contracts

import (
	"context"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/role"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/user"
)

type UserRepository interface {
	Create(
		ctx context.Context,
		user user.User,
		role role.Name,
	) (user.User, error)

	FindByEmail(
		ctx context.Context,
		email string,
	) (user.User, error)
}
