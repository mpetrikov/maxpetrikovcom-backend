package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/labsession"
)

type LabSessionRepository interface {
	Create(
		ctx context.Context,
		session labsession.Session,
	) (labsession.Session, error)

	FindByID(
		ctx context.Context,
		id uuid.UUID,
		userID uuid.UUID,
	) (labsession.Session, error)

	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (labsession.Session, error)

	ListByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]labsession.Session, error)

	Stop(
		ctx context.Context,
		id uuid.UUID,
		userID uuid.UUID,
	) error

	MarkProvisioning(
		ctx context.Context,
		id uuid.UUID,
	) error

	MarkRunning(
		ctx context.Context,
		id uuid.UUID,
	) error
}
