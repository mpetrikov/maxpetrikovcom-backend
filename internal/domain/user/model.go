package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/role"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash *string
	RoleID       int16
	Role         role.Name
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
