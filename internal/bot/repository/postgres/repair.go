package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

func (pg *Pg) CreateRepair(ctx context.Context, repair *entities.Repair) error {
	query := "INSERT INTO repair (producer_name, model_name, name, description, price) VALUES ($1, $2, $3, $4, $5)"
	id, err := pg.db.ExecContextWithReturnID(ctx, query, repair.ProducerName, repair.ModelName, repair.Name, repair.Description, repair.Price)

	if err != nil {
		return err
	}

	repair.ID = uuid.MustParse(id.(string))
	return nil
}

func (pg *Pg) EditRepair(ctx context.Context, repair *entities.Repair) error {
	query := "UPDATE repair SET producer_name = $1, model_name = $2, name = $3, description = $4, price = $5 WHERE id = $6"
	_, err := pg.db.ExecContext(ctx, query, repair.ProducerName, repair.ModelName, repair.Name, repair.Description, repair.Price, repair.ID)

	return err
}

func (pg *Pg) GetProducersRepair(ctx context.Context) ([]string, error) {
	var categories []string

	query := "SELECT DISTINCT(repair.producer_name) FROM repair"
	err := pg.db.SelectContext(ctx, &categories, query)

	if err != nil {
		return nil, err
	} else if len(categories) < 1 {
		return nil, sql.ErrNoRows
	}

	return categories, nil
}

func (pg *Pg) GetModelsRepair(ctx context.Context, producerName string) ([]string, error) {
	var models []string

	query := "SELECT DISTINCT(repair.model_name) FROM repair WHERE lower(producer_name) = lower($1)"
	err := pg.db.SelectContext(ctx, &models, query, producerName)

	if err != nil {
		return nil, err
	} else if len(models) < 1 {
		return nil, sql.ErrNoRows
	}

	return models, nil
}

func (pg *Pg) GetRepairs(ctx context.Context, modelName string) ([]entities.Repair, error) {
	var repairs []entities.Repair

	query := "SELECT * FROM repair WHERE lower(model_name) = lower($1)"
	err := pg.db.SelectContext(ctx, &repairs, query, modelName)

	if err != nil {
		return nil, err
	} else if len(repairs) < 1 {
		return nil, sql.ErrNoRows
	}

	return repairs, nil
}

func (pg *Pg) GetAllRepairs(ctx context.Context) ([]entities.Repair, error) {
	var repairs []entities.Repair

	query := "SELECT * FROM repair"
	err := pg.db.SelectContext(ctx, &repairs, query)

	return repairs, err
}

func (pg *Pg) GetRepairByModelAndName(ctx context.Context, modelName string, repairName string) (*entities.Repair, error) {
	var repair entities.Repair

	query := "SELECT * FROM repair WHERE lower(model_name) = lower($1) AND lower(name) = lower($2)"
	err := pg.db.GetContext(ctx, &repair, query, modelName, repairName)

	return &repair, err
}

func (pg *Pg) GetRepairByModelAndID(ctx context.Context, modelName string, id uuid.UUID) (*entities.Repair, error) {
	var repair entities.Repair

	query := "SELECT * FROM repair WHERE lower(model_name) = lower($1) AND id = $2"
	err := pg.db.GetContext(ctx, &repair, query, modelName, id)

	return &repair, err
}

func (pg *Pg) GetRepairByID(ctx context.Context, id uuid.UUID) (*entities.Repair, error) {
	var repair entities.Repair

	query := "SELECT * FROM repair WHERE id = $1"
	err := pg.db.GetContext(ctx, &repair, query, id)

	return &repair, err
}

func (pg *Pg) DeleteRepairByID(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM repair WHERE id = $1"
	_, err := pg.db.ExecContext(ctx, query, id)

	return err
}
