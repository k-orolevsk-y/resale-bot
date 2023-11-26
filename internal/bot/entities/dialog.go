package entities

import (
	"time"

	"github.com/google/uuid"
)

type Dialog struct {
	ID        uuid.UUID  `db:"id"`
	UserID    int64      `db:"user_id"`
	ManagerID int64      `db:"manager_id"`
	StartedAt time.Time  `db:"started_at"`
	EndedAt   *time.Time `db:"ended_at"`
}
