package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

func (pg *Pg) GetRepairs(ctx context.Context, modelName string) ([]entities.Repair, error) {
	var repairs []entities.Repair

	query := "SELECT repair.* FROM models_repair JOIN repair ON models_repair.id = repair.model_id WHERE lower(models_repair.name) = lower($1)"
	err := pg.db.SelectContext(ctx, &repairs, query, modelName)

	if err != nil {
		return nil, err
	} else if len(repairs) < 1 {
		return nil, sql.ErrNoRows
	}

	return repairs, nil
}

func (pg *Pg) GetRepairWithModelAndCategory(ctx context.Context, modelID uuid.UUID, repairName string) (*entities.RepairWithModelAndCategory, error) {
	var repair entities.RepairWithModelAndCategory

	query := `
		SELECT
			repair.*, models_repair.name AS model_name, categories_repair.name as category_name
		FROM models_repair
			JOIN repair
				ON models_repair.id = repair.model_id
			JOIN categories_repair
				ON categories_repair.id = models_repair.category_repair_id
		WHERE
			models_repair.id = $1
		  AND
			lower(repair.name) = lower($2)
	`
	err := pg.db.GetContext(ctx, &repair, query, modelID, repairName)

	return &repair, err
}

func (pg *Pg) GetRepairWithModelAndCategoryByID(ctx context.Context, modelID, repairID uuid.UUID) (*entities.RepairWithModelAndCategory, error) {
	var repair entities.RepairWithModelAndCategory

	query := `
		SELECT
			repair.*, models_repair.name AS model_name, categories_repair.name as category_name
		FROM models_repair
			JOIN repair
				ON models_repair.id = repair.model_id
			JOIN categories_repair
				ON categories_repair.id = models_repair.category_repair_id
		WHERE
			models_repair.id = $1
		  AND
			repair.id = $2
	`
	err := pg.db.GetContext(ctx, &repair, query, modelID, repairID)

	return &repair, err
}
