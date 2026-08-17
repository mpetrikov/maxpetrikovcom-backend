package job

import "github.com/google/uuid"

type LabCreate struct {
	LabSessionID uuid.UUID `json:"lab_session_id"`
	LabID        uuid.UUID `json:"lab_id"`
	UserID       uuid.UUID `json:"user_id"`
}
