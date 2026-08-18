package lab

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/lab"
	domainlab "github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/lab"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/contracts"
)

type Service struct {
	labRepository contracts.LabRepository
}

type CreateInput struct {
	Slug           string
	Title          string
	Description    string
	Difficulty     domainlab.Difficulty
	TimeoutMinutes int
	Image          string
	CPULimit       string
	MemoryLimit    string
	IsPublished    bool
}

func NewService(
	labRepository contracts.LabRepository,
) *Service {
	return &Service{
		labRepository: labRepository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	input CreateInput,
) (domainlab.Lab, error) {
	newLab := domainlab.Lab{
		ID:             uuid.New(),
		Slug:           strings.TrimSpace(strings.ToLower(input.Slug)),
		Title:          strings.TrimSpace(input.Title),
		Description:    input.Description,
		Difficulty:     input.Difficulty,
		TimeoutMinutes: input.TimeoutMinutes,
		Image:          strings.TrimSpace(input.Image),
		CPULimit:       strings.TrimSpace(input.CPULimit),
		MemoryLimit:    strings.TrimSpace(input.MemoryLimit),
		IsPublished:    input.IsPublished,
	}

	return s.labRepository.Create(ctx, newLab)
}

func (s *Service) FindBySlug(
	ctx context.Context,
	slug string,
) (domainlab.Lab, error) {
	return s.labRepository.FindBySlug(ctx, slug)
}

func (s *Service) ListPublished(
	ctx context.Context,
) ([]domainlab.Lab, error) {
	return s.labRepository.ListPublished(ctx)
}

func (s *Service) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (lab.Lab, error) {
	return s.labRepository.FindByID(ctx, id)
}
