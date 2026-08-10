package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/user"
)

type TokenService struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

type AccessTokenClaims struct {
	RoleID int16 `json:"role_id"`

	jwt.RegisteredClaims
}

func NewTokenService(
	secret string,
	issuer string,
	ttl time.Duration,
) *TokenService {
	return &TokenService{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
	}
}

func (s *TokenService) GenerateAccessToken(
	currentUser user.User,
) (string, error) {
	now := time.Now()

	claims := AccessTokenClaims{
		RoleID: currentUser.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   currentUser.ID.String(),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signedToken, nil
}
