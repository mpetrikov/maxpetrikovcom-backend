package labsession

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending      Status = "pending"
	StatusProvisioning Status = "provisioning"
	StatusRunning      Status = "running"
	StatusStopping     Status = "stopping"
	StatusStopped      Status = "stopped"
	StatusExpired      Status = "expired"
	StatusFailed       Status = "failed"
)

type Session struct {
	ID     uuid.UUID
	LabID  uuid.UUID
	UserID uuid.UUID

	Status Status

	Namespace *string
	PodName   *string

	CreatedAt  time.Time
	StartedAt  *time.Time
	ExpiresAt  time.Time
	FinishedAt *time.Time

	FailureReason *string
}
