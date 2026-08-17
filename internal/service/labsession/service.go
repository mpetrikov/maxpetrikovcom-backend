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
	labs        repositorycontracts.LabRepository
	labSessions repositorycontracts.LabSessionRepository
	publisher   queuecontracts.LabJobPublisher
}

func NewService(
	labs repositorycontracts.LabRepository,
	labSessions repositorycontracts.LabSessionRepository,
	publisher queuecontracts.LabJobPublisher,
) *Service {
	return &Service{
		labs:        labs,
		labSessions: labSessions,
		publisher:   publisher,
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

	created, err := s.labSessions.Create(ctx, session)
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
	return s.labSessions.FindByID(
		ctx,
		id,
		userID,
	)
}

func (s *Service) ListForUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]labsession.Session, error) {
	return s.labSessions.ListByUserID(
		ctx,
		userID,
	)
}

func (s *Service) Stop(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) error {
	return s.labSessions.Stop(
		ctx,
		id,
		userID,
	)
}

func (s *Service) MarkProvisioning(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.labSessions.MarkProvisioning(ctx, id)
}
