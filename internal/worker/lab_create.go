package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/job"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/labsession"
)

func (w *Worker) handleLabCreate(
	ctx context.Context,
	body []byte,
) error {
	var message job.LabCreate

	if err := json.Unmarshal(body, &message); err != nil {
		return fmt.Errorf(
			"decode lab.create message: %w",
			err,
		)
	}

	w.logger.Info(
		"received lab.create",
		"lab_session_id", message.LabSessionID,
		"lab_id", message.LabID,
		"user_id", message.UserID,
	)

	if err := w.labSessionService.MarkProvisioning(
		ctx,
		message.LabSessionID,
	); err != nil {
		if errors.Is(err, labsession.ErrNotFound) {
			return fmt.Errorf(
				"lab session not pending or not found: %w",
				err,
			)
		}

		return err
	}

	w.logger.Info(
		"lab session moved to provisioning",
		"lab_session_id", message.LabSessionID,
	)

	return nil
}
