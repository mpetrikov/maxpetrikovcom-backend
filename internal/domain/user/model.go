package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash *string
	RoleID       int16
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
