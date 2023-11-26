package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Tag          string    `db:"tag"`
	IsManager    bool      `db:"is_manager"`
	RegisteredAt time.Time `db:"registered_at"`
}
