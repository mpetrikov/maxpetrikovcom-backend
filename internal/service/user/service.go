package user

import (
	"context"

	"github.com/google/uuid"

	domainuser "github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/user"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/repository/contracts"
)

type Service struct {
	users contracts.UserRepository
}

func NewService(
	users contracts.UserRepository,
) *Service {
	return &Service{
		users: users,
	}
}

func (s *Service) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domainuser.User, error) {
	return s.users.FindByID(ctx, id)
}
