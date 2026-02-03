package domain_pet

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	Kind           string     `json:"kind"`
	Breed          string     `json:"breed"`
	AgeDays        int        `json:"age_days"`
	Level          int        `json:"level"`
	LastTrainingAt *time.Time `json:"last_training_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}
