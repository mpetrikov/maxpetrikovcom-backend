package labexecution

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/job"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/lab"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/labsession"
	runnercontracts "github.com/maxpetrikov/maxpetrikovcom-backend/internal/runner/contracts"
	labservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/lab"
	labsessionservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/labsession"
)

type Service struct {
	labService        *labservice.Service
	labSessionService *labsessionservice.Service
	labRunner         runnercontracts.LabRunner
	logger            *slog.Logger
}

const expiredSessionBatchSize = 100

func NewService(
	labService *labservice.Service,
	labSessionService *labsessionservice.Service,
	labRunner runnercontracts.LabRunner,
	logger *slog.Logger,
) *Service {
	return &Service{
		labService:        labService,
		labSessionService: labSessionService,
		labRunner:         labRunner,
		logger:            logger,
	}
}

func (s *Service) Start(ctx context.Context, message job.LabCreate) error {
	labSession, ok, err := s.getStartableSession(
		ctx,
		message.LabSessionID,
	)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	if isExpired(labSession, time.Now()) {
		if err := s.expireSession(ctx, labSession); err != nil {
			return fmt.Errorf(
				"expire lab session before start: %w",
				err,
			)
		}

		s.logger.Info(
			"lab session expired before runtime start",
			"lab_session_id", message.LabSessionID,
		)

		return nil
	}

	currentLab, err := s.labService.FindByID(
		ctx,
		message.LabID,
	)
	if err != nil {
		return fmt.Errorf("find lab: %w", err)
	}

	startResult, err := s.startRuntime(
		ctx,
		labSession,
		currentLab,
		message.LabSessionID,
	)
	if err != nil {
		return err
	}

	return s.markRunning(
		ctx,
		message.LabSessionID,
		startResult,
	)
}

func (s *Service) isLabCreateAlreadySettled(
	ctx context.Context,
	labSessionID uuid.UUID,
) bool {
	labSession, err := s.labSessionService.GetById(
		ctx,
		labSessionID,
	)
	if err != nil {
		return false
	}

	switch labSession.Status {
	case labsession.StatusRunning,
		labsession.StatusStopping,
		labsession.StatusStopped,
		labsession.StatusExpired,
		labsession.StatusFailed:
		s.logger.Info(
			"lab.create skipped for already settled lab session",
			"lab_session_id", labSessionID,
			"status", labSession.Status,
		)

		return true

	default:
		return false
	}
}

func (s *Service) getStartableSession(
	ctx context.Context,
	labSessionID uuid.UUID,
) (labsession.Session, bool, error) {
	if err := s.labSessionService.MarkProvisioning(
		ctx,
		labSessionID,
	); err != nil {
		if errors.Is(err, labsession.ErrNotFound) {
			if s.isLabCreateAlreadySettled(ctx, labSessionID) {
				return labsession.Session{}, false, nil
			}

			return labsession.Session{}, false, fmt.Errorf(
				"lab session not pending or not found: %w",
				err,
			)
		}

		return labsession.Session{}, false, err
	}

	labSession, err := s.labSessionService.GetById(
		ctx,
		labSessionID,
	)
	if err != nil {
		if errors.Is(err, labsession.ErrNotFound) {
			return labsession.Session{}, false, fmt.Errorf(
				"lab session not pending or not found: %w",
				err,
			)
		}

		return labsession.Session{}, false, err
	}

	return labSession, true, nil
}

func (s *Service) startRuntime(
	ctx context.Context,
	labSession labsession.Session,
	currentLab lab.Lab,
	labSessionID uuid.UUID,
) (runnercontracts.StartResult, error) {
	startResult, err := s.labRunner.Start(
		ctx,
		labSession,
		currentLab,
	)
	if err == nil {
		return startResult, nil
	}

	startErr := fmt.Errorf("start lab environment: %w", err)

	if stopErr := s.labRunner.Stop(
		ctx,
		labSession,
	); stopErr != nil {
		return runnercontracts.StartResult{}, fmt.Errorf(
			"cleanup failed lab environment after start error: %w: %v",
			stopErr,
			startErr,
		)
	}

	if markErr := s.labSessionService.MarkFailed(
		ctx,
		labSessionID,
		startErr.Error(),
	); markErr != nil {
		return runnercontracts.StartResult{}, fmt.Errorf(
			"mark lab session failed after start error: %w: %v",
			markErr,
			startErr,
		)
	}

	return runnercontracts.StartResult{}, NewTerminalError(startErr)
}

func (s *Service) markRunning(
	ctx context.Context,
	labSessionID uuid.UUID,
	startResult runnercontracts.StartResult,
) error {
	if err := s.labSessionService.MarkRunning(
		ctx,
		labSessionID,
		startResult.Namespace,
		startResult.PodName,
	); err != nil {
		return fmt.Errorf("mark lab session running: %w", err)
	}

	s.logger.Info(
		"lab session is running",
		"lab_session_id", labSessionID,
	)

	return nil
}

func (s *Service) ExpireSessions(ctx context.Context) error {
	sessions, err := s.labSessionService.ListExpiredActive(
		ctx,
		time.Now(),
		expiredSessionBatchSize,
	)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if err := s.expireSession(ctx, session); err != nil {
			return fmt.Errorf(
				"expire lab session %s: %w",
				session.ID,
				err,
			)
		}

		s.logger.Info(
			"lab session expired",
			"lab_session_id", session.ID,
		)
	}

	return nil
}

func (s *Service) StopRuntime(ctx context.Context, message job.LabStop) error {
	labSession, err := s.labSessionService.GetById(
		ctx,
		message.LabSessionID,
	)
	if err != nil {
		if errors.Is(err, labsession.ErrNotFound) {
			return fmt.Errorf(
				"lab session not found: %w",
				err,
			)
		}

		return err
	}

	if err := s.labRunner.Stop(
		ctx,
		labSession,
	); err != nil {
		return fmt.Errorf("stop lab environment: %w", err)
	}

	if err := s.labSessionService.MarkStopped(
		ctx,
		message.LabSessionID,
	); err != nil {
		return fmt.Errorf("mark lab session stopped: %w", err)
	}

	s.logger.Info(
		"lab session is stopped",
		"lab_session_id", message.LabSessionID,
	)

	return nil
}

func (s *Service) expireSession(
	ctx context.Context,
	labSession labsession.Session,
) error {
	if labSession.Status != labsession.StatusPending {
		if err := s.labRunner.Stop(
			ctx,
			labSession,
		); err != nil {
			return fmt.Errorf("stop expired lab environment: %w", err)
		}
	}

	if err := s.labSessionService.MarkExpired(
		ctx,
		labSession.ID,
	); err != nil {
		if errors.Is(err, labsession.ErrNotFound) {
			return nil
		}

		return fmt.Errorf("mark lab session expired: %w", err)
	}

	return nil
}

func isExpired(
	labSession labsession.Session,
	now time.Time,
) bool {
	return !labSession.ExpiresAt.After(now)
}
