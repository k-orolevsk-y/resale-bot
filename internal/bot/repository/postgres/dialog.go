package postgres

import (
	"context"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

func (pg *Pg) CreateDialog(ctx context.Context, dialog *entities.Dialog) error {
	query := "INSERT INTO dialog (user_id, manager_id) VALUES ($1, $2)"
	id, err := pg.db.ExecContextWithReturnID(ctx, query, dialog.UserID, dialog.ManagerID)

	if err != nil {
		return err
	}

	dialog.ID = uuid.MustParse(id.(string))
	return nil
}

func (pg *Pg) EditDialog(ctx context.Context, dialog *entities.Dialog) error {
	query := "UPDATE dialog SET user_id = $1, manager_id = $2, ended_at = $3 WHERE id = $4"
	_, err := pg.db.ExecContext(ctx, query, dialog.UserID, dialog.ManagerID, dialog.EndedAt, dialog.ID)

	return err
}

func (pg *Pg) GetDialogByTalkerID(ctx context.Context, talkerID int64) (*entities.Dialog, error) {
	var dialog entities.Dialog

	query := "SELECT * FROM dialog WHERE (user_id = $1 OR manager_id = $2) AND ended_at IS NULL"
	err := pg.db.GetContext(ctx, &dialog, query, talkerID, talkerID)

	return &dialog, err
}
