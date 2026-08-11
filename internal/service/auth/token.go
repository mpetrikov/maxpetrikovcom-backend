package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/role"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/user"
)

type TokenService struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

type AccessTokenClaims struct {
	Role role.Name `json:"role"`

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
		Role: currentUser.Role,
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

func (s *TokenService) ParseAccessToken(
	tokenString string,
) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessTokenClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf(
					"unexpected signing method: %s",
					token.Method.Alg(),
				)
			}

			return s.secret, nil
		},
		jwt.WithIssuer(s.issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("parse access token: %w", err)
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}

	return claims, nil
}
