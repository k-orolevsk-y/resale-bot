package postgres

import (
	"context"
	"database/sql"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

func (pg *Pg) GetModelsRepair(ctx context.Context, categoryName string) ([]entities.ModelRepair, error) {
	var models []entities.ModelRepair

	query := "SELECT models_repair.* FROM categories_repair JOIN models_repair ON models_repair.category_repair_id = categories_repair.id WHERE lower(categories_repair.name) = lower($1)"
	err := pg.db.SelectContext(ctx, &models, query, categoryName)

	if err != nil {
		return nil, err
	} else if len(models) < 1 {
		return nil, sql.ErrNoRows
	}

	return models, nil
}
