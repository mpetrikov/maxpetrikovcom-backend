package worker

import "context"

func (w *Worker) handleLabSessionExpiration(
	ctx context.Context,
) error {
	return w.labExecutionService.ExpireSessions(ctx)
}
