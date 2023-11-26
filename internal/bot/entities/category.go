package entities

import "github.com/google/uuid"

type Category struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
	Type int       `db:"c_type"`
}
