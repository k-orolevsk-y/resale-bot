package postgres

import (
	"context"
	"database/sql"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

func (pg *Pg) GetCategoriesRepair(ctx context.Context) ([]entities.CategoryRepair, error) {
	var categories []entities.CategoryRepair

	query := "SELECT * FROM categories_repair"
	err := pg.db.SelectContext(ctx, &categories, query)

	if err != nil {
		return nil, err
	} else if len(categories) < 1 {
		return nil, sql.ErrNoRows
	}

	return categories, nil
}
