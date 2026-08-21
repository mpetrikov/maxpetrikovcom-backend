package lab

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	domainlab "github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/lab"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/contracts"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
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
	image := strings.TrimSpace(input.Image)
	cpuLimit := strings.TrimSpace(input.CPULimit)
	memoryLimit := strings.TrimSpace(input.MemoryLimit)

	if err := validateRuntimeConfig(
		image,
		cpuLimit,
		memoryLimit,
	); err != nil {
		return domainlab.Lab{}, err
	}

	newLab := domainlab.Lab{
		ID:             uuid.New(),
		Slug:           strings.TrimSpace(strings.ToLower(input.Slug)),
		Title:          strings.TrimSpace(input.Title),
		Description:    input.Description,
		Difficulty:     input.Difficulty,
		TimeoutMinutes: input.TimeoutMinutes,
		Image:          image,
		CPULimit:       cpuLimit,
		MemoryLimit:    memoryLimit,
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
) (domainlab.Lab, error) {
	return s.labRepository.FindByID(ctx, id)
}

func validateRuntimeConfig(
	image string,
	cpuLimit string,
	memoryLimit string,
) error {
	if image == "" {
		return fmt.Errorf(
			"%w: image is required",
			domainlab.ErrInvalidInput,
		)
	}

	if _, err := k8sresource.ParseQuantity(cpuLimit); err != nil {
		return fmt.Errorf(
			"%w: invalid CPU limit: %v",
			domainlab.ErrInvalidInput,
			err,
		)
	}

	if _, err := k8sresource.ParseQuantity(memoryLimit); err != nil {
		return fmt.Errorf(
			"%w: invalid memory limit: %v",
			domainlab.ErrInvalidInput,
			err,
		)
	}

	return nil
}
