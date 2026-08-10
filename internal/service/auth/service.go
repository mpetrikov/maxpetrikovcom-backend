package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/contracts"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/role"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/user"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/useridentity"
)

type Service struct {
	users      contracts.UserRepository
	identities contracts.UserIdentityRepository
}

func NewService(users contracts.UserRepository, identities contracts.UserIdentityRepository) *Service {
	return &Service{
		users:      users,
		identities: identities,
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

func (s *Service) LoginWithGoogle(
	ctx context.Context,
	googleUser GoogleUser,
) (user.User, error) {
	existingUser, err := s.identities.FindUserByProvider(
		ctx,
		useridentity.Google,
		googleUser.ID,
	)

	if err == nil {
		return existingUser, nil
	}

	if !errors.Is(err, useridentity.ErrNotFound) {
		return user.User{}, fmt.Errorf(
			"find google identity: %w",
			err,
		)
	}

	if !googleUser.EmailVerified {
		return user.User{}, errors.New("google email is not verified")
	}

	newUser := user.User{
		ID:    uuid.New(),
		Email: strings.ToLower(strings.TrimSpace(googleUser.Email)),
	}

	email := newUser.Email

	newIdentity := useridentity.Identity{
		ID:             uuid.New(),
		UserID:         newUser.ID,
		Provider:       useridentity.Google,
		ProviderUserID: googleUser.ID,
		Email:          &email,
	}

	created, err := s.identities.CreateUserWithIdentity(
		ctx,
		newUser,
		newIdentity,
		role.Student,
	)
	if err != nil {
		return user.User{}, fmt.Errorf(
			"create google user: %w",
			err,
		)
	}

	return created, nil
}
