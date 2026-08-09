package useridentity

import (
	"time"

	"github.com/google/uuid"
)

type Provider string

const (
	Google Provider = "google"
	GitHub Provider = "github"
)

type Identity struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Provider       Provider
	ProviderUserID string
	Email          *string
	CreatedAt      time.Time
}
