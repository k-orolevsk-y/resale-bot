package callback

import (
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

type Repository interface {
}

type service struct {
	rep    Repository
	logger *zap.Logger
}

func ConfigureKeyboardCallbackService(app *app.App) {
	s := &service{
		rep:    app.GetRepository(),
		logger: app.GetLogger(),
	}
	engine := app.GetEngine()

	engine.Callback("cancel_manager", s.CancelManager)
	engine.Group("manager_", func(group bot.Router) {
		group.Use(s.ManagerAccess)

		group.Callback("manager_dialog_start", s.ManagerDialogStart)
	})
}
