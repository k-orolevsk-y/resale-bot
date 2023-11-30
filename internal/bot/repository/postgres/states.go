package postgres

import (
	"context"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

func (pg *Pg) CreateState(ctx context.Context, state *entities.State) error {
	if err := state.EncodeData(); err != nil {
		return err
	}

	query := "INSERT INTO states (id, s_type, data) VALUES ($1, $2, $3)"
	_, err := pg.db.ExecContext(ctx, query, state.ID, state.Type, state.Data)

	return err
}

func (pg *Pg) GetState(ctx context.Context, id string, sType int) (*entities.State, error) {
	var state entities.State

	query := "SELECT * FROM states WHERE id = $1 AND s_type = $2"
	err := pg.db.GetContext(ctx, &state, query, id, sType)

	if err != nil {
		return nil, err
	}

	if err = state.DecodeData(); err != nil {
		return nil, err
	}

	return &state, err
}

func (pg *Pg) DeleteState(ctx context.Context, id string, sType int) error {
	query := "DELETE FROM states WHERE id = $1 AND s_type = $2"
	_, err := pg.db.ExecContext(ctx, query, id, sType)

	return err
}
