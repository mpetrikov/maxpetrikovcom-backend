package refreshsession

import "errors"

var (
	ErrNotFound = errors.New("refresh session not found")
	ErrExpired  = errors.New("refresh session expired")
	ErrRevoked  = errors.New("refresh session revoked")
)
