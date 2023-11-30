package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

func (pg *Pg) CreateReservation(ctx context.Context, reservation *entities.Reservation) error {
	query := "INSERT INTO reservation (user_id, product_id) VALUES ($1, $2)"

	id, err := pg.db.ExecContextWithReturnID(ctx, query, reservation.UserID, reservation.ProductID)
	if err != nil {
		return err
	}

	reservation.ID = uuid.MustParse(id.(string))
	return nil
}

func (pg *Pg) ExistsReservationByProductID(ctx context.Context, productID uuid.UUID) (bool, error) {
	var exists bool

	query := "SELECT true AS exists FROM reservation WHERE product_id = $1 AND completed = 0 LIMIT 1"
	err := pg.db.GetContext(ctx, &exists, query, productID)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	return exists, err
}
