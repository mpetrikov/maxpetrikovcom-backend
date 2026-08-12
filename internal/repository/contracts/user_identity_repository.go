package contracts

import (
	"context"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/role"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/user"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/useridentity"
)

type UserIdentityRepository interface {
	FindUserByProvider(
		ctx context.Context,
		provider useridentity.Provider,
		providerUserID string,
	) (user.User, error)

	CreateUserWithIdentity(
		ctx context.Context,
		user user.User,
		identity useridentity.Identity,
		role role.Name,
	) (user.User, error)
}
