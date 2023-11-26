package entities

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	ProductID uuid.UUID `db:"product_id"`
	CreatedAt time.Time `db:"created_at"`
}
