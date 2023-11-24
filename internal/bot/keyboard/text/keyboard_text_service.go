package text

import (
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
)

type Repository interface {
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
}
