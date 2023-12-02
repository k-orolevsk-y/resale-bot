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

func (pg *Pg) EditReservation(ctx context.Context, reservation *entities.Reservation) error {
	query := "UPDATE reservation SET completed = $1 WHERE id = $2"
	_, err := pg.db.ExecContext(ctx, query, reservation.Completed, reservation.ID)

	return err
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

func (pg *Pg) GetReservationsByUserID(ctx context.Context, userID int64) ([]entities.ReservationWithAdditionalData, error) {
	var reservations []entities.ReservationWithAdditionalData

	query := `
		SELECT
			reservation.*,
			categories.name as category_name,
			CONCAT(products.model, ' - ', products.additional) AS product_full_name
		FROM reservation
			JOIN products ON reservation.product_id = products.id
			JOIN categories ON products.category_id = categories.id
		WHERE user_id = $1
	`
	err := pg.db.SelectContext(ctx, &reservations, query, userID)

	if err != nil {
		return nil, err
	} else if len(reservations) < 1 {
		return nil, sql.ErrNoRows
	}

	return reservations, nil
}

func (pg *Pg) GetReservations(ctx context.Context) ([]entities.ReservationWithAdditionalData, error) {
	var reservations []entities.ReservationWithAdditionalData

	query := `
		SELECT
			reservation.*,
			categories.name as category_name,
			CONCAT(products.model, ' - ', products.additional) AS product_full_name
		FROM reservation
			JOIN products ON reservation.product_id = products.id
			JOIN categories ON products.category_id = categories.id
		ORDER BY reservation.completed <> 0
	`
	err := pg.db.SelectContext(ctx, &reservations, query)

	if err != nil {
		return nil, err
	} else if len(reservations) < 1 {
		return nil, sql.ErrNoRows
	}

	return reservations, nil
}

func (pg *Pg) GetReservationByID(ctx context.Context, id uuid.UUID) (*entities.ReservationWithAdditionalData, error) {
	var reservation entities.ReservationWithAdditionalData

	query := `
		SELECT
			reservation.*,
			categories.name as category_name,
			CONCAT(products.model, ' - ', products.additional) AS product_full_name
		FROM reservation
			JOIN products ON reservation.product_id = products.id
			JOIN categories ON products.category_id = categories.id
		WHERE reservation.id = $1
	`
	err := pg.db.GetContext(ctx, &reservation, query, id)

	return &reservation, err
}

func (pg *Pg) DeleteReservationsByCategoryID(ctx context.Context, categoryID uuid.UUID) error {
	query := "DELETE FROM reservation WHERE product_id IN (SELECT id FROM products WHERE category_id = $1)"
	_, err := pg.db.ExecContext(ctx, query, categoryID)

	return err
}

func (pg *Pg) DeleteReservationsByProductID(ctx context.Context, productID uuid.UUID) error {
	query := "DELETE FROM reservation WHERE product_id = $1"
	_, err := pg.db.ExecContext(ctx, query, productID)

	return err
}
