package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/contracts"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/role"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/user"
)

type Service struct {
	users contracts.UserRepository
}

func NewService(users contracts.UserRepository) *Service {
	return &Service{
		users: users,
	}
}

func (s *Service) Register(
	ctx context.Context,
	email string,
	password string,
) (user.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return user.User{}, fmt.Errorf("hash password: %w", err)
	}

	passwordHash := string(hash)

	newUser := user.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: &passwordHash,
	}

	created, err := s.users.Create(
		ctx,
		newUser,
		role.Student,
	)
	if err != nil {
		return user.User{}, fmt.Errorf("register user: %w", err)
	}

	return created, nil
}
