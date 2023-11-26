package text

import (
	"context"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

type Repository interface {
	CreateDialog(context.Context, *entities.Dialog) error
	EditDialog(ctx context.Context, dialog *entities.Dialog) error
	GetDialogByTalkerID(context.Context, int64) (*entities.Dialog, error)
}

type service struct {
	rep    Repository
	logger *zap.Logger
}

func ConfigureKeyboardTextService(app *app.App) {
	s := &service{
		rep:    app.GetRepository(),
		logger: app.GetLogger(),
	}
	engine := app.GetEngine()

	engine.Message("Главное меню", s.HomeMenu)
	engine.Message("Связаться с менеджером", s.Manager)
	engine.Message("Завершить диалог", s.ExitFromDialog)
}
