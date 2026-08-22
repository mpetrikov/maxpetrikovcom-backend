package contracts

import (
	"context"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/job"
)

type LabJobPublisher interface {
	PublishCreate(
		ctx context.Context,
		job job.LabCreate,
	) error

	PublishStop(
		ctx context.Context,
		job job.LabStop,
	) error
}
