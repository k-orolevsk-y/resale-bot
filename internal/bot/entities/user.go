package entities

import (
	"time"
)

type User struct {
	ID           int64     `db:"id"`
	Tag          string    `db:"tag"`
	IsManager    bool      `db:"is_manager"`
	IsBanned     bool      `db:"is_banned"`
	RegisteredAt time.Time `db:"registered_at"`
}
