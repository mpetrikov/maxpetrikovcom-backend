package contracts

import (
	"context"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/useridentity"
)

type UserIdentityRepository interface {
	Create(
		ctx context.Context,
		identity useridentity.Identity,
	) (useridentity.Identity, error)

	FindByProvider(
		ctx context.Context,
		provider useridentity.Provider,
		providerUserID string,
	) (useridentity.Identity, error)
}
