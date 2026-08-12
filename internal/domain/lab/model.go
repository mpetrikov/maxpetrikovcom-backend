package lab

import (
	"time"

	"github.com/google/uuid"
)

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type Lab struct {
	ID             uuid.UUID
	Slug           string
	Title          string
	Description    string
	Difficulty     Difficulty
	TimeoutMinutes int
	Image          string
	CPULimit       string
	MemoryLimit    string
	IsPublished    bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
