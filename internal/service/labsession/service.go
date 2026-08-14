package labsession

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/labsession"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/contracts"
)

type Service struct {
	labs        contracts.LabRepository
	labSessions contracts.LabSessionRepository
}

func NewService(
	labs contracts.LabRepository,
	labSessions contracts.LabSessionRepository,
) *Service {
	return &Service{
		labs:        labs,
		labSessions: labSessions,
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

	return s.labSessions.Create(ctx, session)
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
