package entities

import "github.com/google/uuid"

type CategoryRepair struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
