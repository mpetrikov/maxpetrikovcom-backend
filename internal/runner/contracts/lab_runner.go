package contracts

import (
	"context"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/lab"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/labsession"
)

type StartResult struct {
	Namespace string
	PodName   string
}

type LabRunner interface {
	Start(
		ctx context.Context,
		session labsession.Session,
		lab lab.Lab,
	) (StartResult, error)

	Stop(
		ctx context.Context,
		session labsession.Session,
	) error
}
