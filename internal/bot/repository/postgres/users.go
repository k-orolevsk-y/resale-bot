package postgres

import (
	"context"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

func (pg *Pg) CreateUser(ctx context.Context, user *entities.User) error {
	query := "INSERT INTO users (id, tag, is_manager) VALUES ($1, $2, $3)"
	_, err := pg.db.ExecContext(ctx, query, user.ID, user.Tag, user.IsManager)

	return err
}

func (pg *Pg) EditUser(ctx context.Context, user *entities.User) error {
	query := "UPDATE users SET tag = $1, is_manager = $2 WHERE id = $3"
	_, err := pg.db.ExecContext(ctx, query, user.Tag, user.IsManager, user.ID)

	return err
}

func (pg *Pg) GetUserByTgID(ctx context.Context, tgID int64) (*entities.User, error) {
	var user entities.User

	query := "SELECT * FROM users WHERE id = $1"
	err := pg.db.GetContext(ctx, &user, query, tgID)

	return &user, err
}
