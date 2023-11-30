package entities

import "github.com/google/uuid"

type ModelRepair struct {
	ID               uuid.UUID `db:"id"`
	CategoryRepairID uuid.UUID `db:"category_repair_id"`
	Name             string    `db:"name"`
}
