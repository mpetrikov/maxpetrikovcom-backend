package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/job"
)

func (w *Worker) handleLabStop(
	ctx context.Context,
	body []byte,
) error {
	var message job.LabStop

	if err := json.Unmarshal(body, &message); err != nil {
		return fmt.Errorf("decode lab.stop message: %w", err)
	}

	w.logger.Info(
		"received lab.stop",
		"lab_session_id", message.LabSessionID,
		"user_id", message.UserID,
	)

	if err := w.labExecutionService.StopRuntime(
		ctx,
		message,
	); err != nil {
		return fmt.Errorf("stop lab execution: %w", err)
	}

	return nil
}
