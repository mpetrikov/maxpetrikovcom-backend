package contracts

import (
	"context"

	"github.com/google/uuid"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/lab"
)

type LabRepository interface {
	Create(
		ctx context.Context,
		lab lab.Lab,
	) (lab.Lab, error)

	FindBySlug(
		ctx context.Context,
		slug string,
	) (lab.Lab, error)

	FindByID(
		ctx context.Context,
		labId uuid.UUID,
	) (lab.Lab, error)

	ListPublished(
		ctx context.Context,
	) ([]lab.Lab, error)
}
