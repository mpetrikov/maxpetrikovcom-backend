package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/contracts"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/refreshsession"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/user"
)

type RefreshTokenService struct {
	sessions contracts.RefreshSessionRepository
	users    contracts.UserRepository
	tokens   *TokenService
	ttl      time.Duration
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func NewRefreshTokenService(
	sessions contracts.RefreshSessionRepository,
	users contracts.UserRepository,
	tokens *TokenService,
	ttl time.Duration,
) *RefreshTokenService {
	return &RefreshTokenService{
		sessions: sessions,
		users:    users,
		tokens:   tokens,
		ttl:      ttl,
	}
}

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *RefreshTokenService) IssueTokenPair(
	ctx context.Context,
	currentUser user.User,
) (TokenPair, error) {
	accessToken, err := s.tokens.GenerateAccessToken(currentUser)
	if err != nil {
		return TokenPair{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}

	session := refreshsession.Session{
		ID:        uuid.New(),
		UserID:    currentUser.ID,
		TokenHash: hashRefreshToken(refreshToken),
		ExpiresAt: time.Now().Add(s.ttl),
	}

	if err := s.sessions.Create(ctx, session); err != nil {
		return TokenPair{}, fmt.Errorf("create refresh session: %w", err)
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *RefreshTokenService) Refresh(
	ctx context.Context,
	refreshToken string,
) (TokenPair, error) {
	tokenHash := hashRefreshToken(refreshToken)

	session, err := s.sessions.FindByTokenHash(
		ctx,
		tokenHash,
	)
	if err != nil {
		return TokenPair{}, err
	}

	if session.RevokedAt != nil {
		return TokenPair{}, refreshsession.ErrRevoked
	}

	if time.Now().After(session.ExpiresAt) {
		return TokenPair{}, refreshsession.ErrExpired
	}

	currentUser, err := s.users.FindByID(
		ctx,
		session.UserID,
	)
	if err != nil {
		return TokenPair{}, fmt.Errorf("find refresh session user: %w", err)
	}

	if err := s.sessions.Revoke(ctx, tokenHash); err != nil {
		return TokenPair{}, fmt.Errorf("revoke refresh session: %w", err)
	}

	return s.IssueTokenPair(ctx, currentUser)
}

func (s *RefreshTokenService) Logout(
	ctx context.Context,
	refreshToken string,
) error {
	tokenHash := hashRefreshToken(refreshToken)

	err := s.sessions.Revoke(ctx, tokenHash)
	if errors.Is(err, refreshsession.ErrNotFound) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("logout refresh session: %w", err)
	}

	return nil
}
