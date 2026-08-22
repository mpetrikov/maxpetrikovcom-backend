package labsession

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/job"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/labsession"
	queuecontracts "github.com/maxpetrikov/maxpetrikovcom-backend/internal/queue/contracts"
	repositorycontracts "github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/contracts"
)

type Service struct {
	labs                 repositorycontracts.LabRepository
	labSessionRepository repositorycontracts.LabSessionRepository
	publisher            queuecontracts.LabJobPublisher
}

const maxFailureReasonLength = 2000

func NewService(
	labs repositorycontracts.LabRepository,
	labSessionRepository repositorycontracts.LabSessionRepository,
	publisher queuecontracts.LabJobPublisher,
) *Service {
	return &Service{
		labs:                 labs,
		labSessionRepository: labSessionRepository,
		publisher:            publisher,
	}
}

func (s *Service) Create(
	ctx context.Context,
	slug string,
	userID uuid.UUID,
) (labsession.Session, error) {
	lab, err := s.labs.FindBySlug(ctx, slug)
	if err != nil {
		return labsession.Session{}, err
	}

	if !lab.IsPublished {
		return labsession.Session{},
			labsession.ErrNotFound
	}

	now := time.Now()

	session := labsession.Session{
		ID:     uuid.New(),
		LabID:  lab.ID,
		UserID: userID,
		Status: labsession.StatusPending,

		ExpiresAt: now.Add(
			time.Duration(lab.TimeoutMinutes) * time.Minute,
		),
	}

	created, err := s.labSessionRepository.Create(ctx, session)
	if err != nil {
		return labsession.Session{}, err
	}

	err = s.publisher.PublishCreate(
		ctx,
		job.LabCreate{
			LabSessionID: created.ID,
			LabID:        created.LabID,
			UserID:       created.UserID,
		},
	)
	if err != nil {
		return labsession.Session{},
			fmt.Errorf("publish lab create job: %w", err)
	}

	return created, nil
}

func (s *Service) Get(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) (labsession.Session, error) {
	return s.labSessionRepository.FindByID(
		ctx,
		id,
		userID,
	)
}

func (s *Service) ListForUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]labsession.Session, error) {
	return s.labSessionRepository.ListByUserID(
		ctx,
		userID,
	)
}

func (s *Service) ListExpiredActive(
	ctx context.Context,
	now time.Time,
	limit int,
) ([]labsession.Session, error) {
	return s.labSessionRepository.ListExpiredActive(
		ctx,
		now,
		limit,
	)
}

func (s *Service) RequestStop(
	ctx context.Context,
	labSessionId uuid.UUID,
	userID uuid.UUID,
) error {
	if err := s.labSessionRepository.MarkStopping(
		ctx,
		labSessionId,
		userID,
	); err != nil {
		return err
	}

	err := s.publisher.PublishStop(
		ctx,
		job.LabStop{
			LabSessionID: labSessionId,
			UserID:       userID,
		},
	)
	if err != nil {
		return fmt.Errorf("publish lab stop job: %w", err)
	}

	return nil
}

func (s *Service) MarkProvisioning(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.labSessionRepository.MarkProvisioning(ctx, id)
}

func (s *Service) MarkRunning(
	ctx context.Context,
	id uuid.UUID,
	namespace string,
	podName string,
) error {
	return s.labSessionRepository.MarkRunning(
		ctx,
		id,
		namespace,
		podName,
	)
}

func (s *Service) GetById(
	ctx context.Context,
	sessionId uuid.UUID,
) (labsession.Session, error) {
	return s.labSessionRepository.GetByID(ctx, sessionId)
}

func (s *Service) MarkStopped(
	ctx context.Context,
	labSessionId uuid.UUID,
) error {
	return s.labSessionRepository.MarkStopped(ctx, labSessionId)
}

func (s *Service) MarkExpired(
	ctx context.Context,
	labSessionId uuid.UUID,
) error {
	return s.labSessionRepository.MarkExpired(ctx, labSessionId)
}

func (s *Service) MarkFailed(
	ctx context.Context,
	labSessionId uuid.UUID,
	reason string,
) error {
	return s.labSessionRepository.MarkFailed(
		ctx,
		labSessionId,
		truncateFailureReason(reason),
	)
}

func truncateFailureReason(reason string) string {
	runes := []rune(reason)
	if len(runes) <= maxFailureReasonLength {
		return reason
	}

	return string(runes[:maxFailureReasonLength])
}
