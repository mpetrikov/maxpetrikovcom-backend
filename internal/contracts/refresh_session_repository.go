package contracts

import (
	"context"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/refreshsession"
)

type RefreshSessionRepository interface {
	Create(
		ctx context.Context,
		session refreshsession.Session,
	) error

	FindByTokenHash(
		ctx context.Context,
		tokenHash string,
	) (refreshsession.Session, error)

	Revoke(
		ctx context.Context,
		tokenHash string,
	) error
}
