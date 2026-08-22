package labexecution

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/job"
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
	if err := s.labSessionService.MarkProvisioning(
		ctx,
		message.LabSessionID,
	); err != nil {
		if errors.Is(err, labsession.ErrNotFound) {
			if s.isLabCreateAlreadySettled(ctx, message.LabSessionID) {
				return nil
			}

			return fmt.Errorf(
				"lab session not pending or not found: %w",
				err,
			)
		}

		return err
	}

	labSession, err := s.labSessionService.GetById(
		ctx,
		message.LabSessionID,
	)
	if err != nil {
		if errors.Is(err, labsession.ErrNotFound) {
			return fmt.Errorf(
				"lab session not pending or not found: %w",
				err,
			)
		}

		return err
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

	startResult, err := s.labRunner.Start(
		ctx,
		labSession,
		currentLab,
	)
	if err != nil {
		return fmt.Errorf("start lab environment: %w", err)
	}

	if err := s.labSessionService.MarkRunning(
		ctx,
		message.LabSessionID,
		startResult.Namespace,
		startResult.PodName,
	); err != nil {
		return fmt.Errorf("mark lab session running: %w", err)
	}

	s.logger.Info(
		"lab session is running",
		"lab_session_id", message.LabSessionID,
	)

	return nil
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
