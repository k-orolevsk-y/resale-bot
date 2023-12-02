package postgres

import (
	"context"
	"fmt"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

func (pg *Pg) CreateUser(ctx context.Context, user *entities.User) error {
	query := "INSERT INTO users (id, tag, is_manager) VALUES ($1, $2, $3)"
	_, err := pg.db.ExecContext(ctx, query, user.ID, user.Tag, user.IsManager)

	return err
}

func (pg *Pg) EditUser(ctx context.Context, user *entities.User) error {
	query := "UPDATE users SET tag = $1, is_manager = $2, is_banned = $3 WHERE id = $4"
	_, err := pg.db.ExecContext(ctx, query, user.Tag, user.IsManager, user.IsBanned, user.ID)

	return err
}

func (pg *Pg) FindUser(ctx context.Context, q interface{}) (*entities.User, error) {
	var user entities.User

	query := "SELECT * FROM users WHERE text(id) = $1 OR replace(tag, '@', '') = replace($1, '@', '')"
	err := pg.db.GetContext(ctx, &user, query, fmt.Sprint(q))

	return &user, err
}

func (pg *Pg) GetUserByTgID(ctx context.Context, tgID int64) (*entities.User, error) {
	var user entities.User

	query := "SELECT * FROM users WHERE id = $1"
	err := pg.db.GetContext(ctx, &user, query, tgID)

	return &user, err
}

func (pg *Pg) GetUserIdsWhoManager(ctx context.Context) ([]int64, error) {
	var userIds []int64

	query := "SELECT id FROM users WHERE is_manager = true"
	err := pg.db.SelectContext(ctx, &userIds, query)

	return userIds, err
}
