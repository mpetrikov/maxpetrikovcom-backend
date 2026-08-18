package fake

import (
	"context"
	"log/slog"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/lab"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/labsession"
)

type LabRunner struct {
	logger *slog.Logger
}

func NewLabRunner(
	logger *slog.Logger,
) *LabRunner {
	return &LabRunner{
		logger: logger,
	}
}

func (r *LabRunner) Start(
	ctx context.Context,
	session labsession.Session,
	lab lab.Lab,
) error {
	r.logger.Info(
		"starting fake lab environment",
		"lab_session_id", session.ID,
		"lab_id", lab.ID,
		"image", lab.Image,
	)

	return nil
}

func (r *LabRunner) Stop(
	ctx context.Context,
	session labsession.Session,
) error {
	r.logger.Info(
		"stopping fake lab environment",
		"lab_session_id", session.ID,
	)

	return nil
}
