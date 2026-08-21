package labexecution

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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
			return fmt.Errorf(
				"lab session not pending or not found: %w",
				err,
			)
		}

		return err
	}

	currentLab, err := s.labService.FindByID(
		ctx,
		message.LabID,
	)
	if err != nil {
		return fmt.Errorf("find lab: %w", err)
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
