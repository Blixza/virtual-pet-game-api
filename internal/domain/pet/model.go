package domain_pet

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID             uuid.UUID
	Name           string
	Kind           string
	Breed          string
	AgeDays        int
	Level          int
	LastTrainingAt sql.NullTime
	CreatedAt      time.Time
	UpdatedAt      sql.NullTime
}
