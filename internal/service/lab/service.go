package lab

import (
	"context"
	"strings"

	"github.com/google/uuid"

	domainlab "github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/lab"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/contracts"
)

type Service struct {
	labs contracts.LabRepository
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
	labs contracts.LabRepository,
) *Service {
	return &Service{
		labs: labs,
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

	return s.labs.Create(ctx, newLab)
}

func (s *Service) FindBySlug(
	ctx context.Context,
	slug string,
) (domainlab.Lab, error) {
	return s.labs.FindBySlug(ctx, slug)
}

func (s *Service) ListPublished(
	ctx context.Context,
) ([]domainlab.Lab, error) {
	return s.labs.ListPublished(ctx)
}
