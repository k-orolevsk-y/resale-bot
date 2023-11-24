package commands

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

func ConfigureCommandsService(app *app.App) {
	s := &service{
		rep:    app.GetRepository(),
		logger: app.GetLogger(),
	}
	engine := app.GetEngine()

	engine.Command("start", s.Start)
}
