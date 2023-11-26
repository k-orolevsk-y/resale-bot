package entities

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Product struct {
	ID              uuid.UUID      `db:"id"`
	CategoryID      uuid.UUID      `db:"category_id"`
	Producer        string         `db:"producer"`
	Model           string         `db:"model"`
	Additional      string         `db:"additional"`
	OperatingSystem int            `db:"operating_system"`
	Description     string         `db:"description"`
	Photos          pq.StringArray `db:"photos"`
}
