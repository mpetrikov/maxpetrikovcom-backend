package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/job"
)

func (w *Worker) handleLabCreate(
	ctx context.Context,
	body []byte,
) error {
	var message job.LabCreate

	if err := json.Unmarshal(body, &message); err != nil {
		return fmt.Errorf("decode lab.create message: %w", err)
	}

	w.logger.Info(
		"received lab.create",
		"lab_session_id", message.LabSessionID,
		"lab_id", message.LabID,
		"user_id", message.UserID,
	)

	if err := w.labExecutionService.Start(
		ctx,
		message,
	); err != nil {
		return fmt.Errorf("start lab execution: %w", err)
	}

	return nil
}
