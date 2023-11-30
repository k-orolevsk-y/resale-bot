package manager

import (
	"context"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

type Repository interface {
	EditDialog(context.Context, *entities.Dialog) error
	GetDialogByTalkerID(context.Context, int64) (*entities.Dialog, error)
}

type keyboardCallbackManagerService struct {
	rep    Repository
	logger *zap.Logger
}

func ConfigureKeyboardCallbackManagerService(app *app.App) {
	service := keyboardCallbackManagerService{
		rep:    app.GetRepository(),
		logger: app.GetLogger(),
	}
	engine := app.GetEngine()

	engine.Group("manager_", func(group bot.Router) {
		group.Use(service.ManagerAccess)

		group.Callback("manager_dialog_start", service.ManagerDialogStart)
	})
}
