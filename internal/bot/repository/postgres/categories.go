package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

func (pg *Pg) CreateCategory(ctx context.Context, category *entities.Category) error {
	query := "INSERT INTO categories (name, c_type) VALUES ($1, $2)"
	id, err := pg.db.ExecContextWithReturnID(ctx, query, category.Name, category.Type)

	if err != nil {
		return err
	}

	category.ID = uuid.MustParse(id.(string))
	return nil
}

func (pg *Pg) GetCategoriesByType(ctx context.Context, cType int) ([]entities.Category, error) {
	var category []entities.Category

	query := "SELECT * FROM categories WHERE c_type = $1"
	err := pg.db.SelectContext(ctx, &category, query, cType)

	if err != nil {
		return nil, err
	} else if len(category) < 1 {
		return nil, sql.ErrNoRows
	}

	return category, err
}

func (pg *Pg) DeleteCategory(ctx context.Context, category *entities.Category) error {
	query := "DELETE FROM categories WHERE id = $1"
	_, err := pg.db.ExecContext(ctx, query, category.ID)

	return err
}
