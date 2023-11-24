package middlewares

import (
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
)

type service struct {
	logger *zap.Logger
}

func ConfigureMiddlewaresService(app *app.App) {
	s := &service{
		logger: app.GetLogger(),
	}
	engine := app.GetEngine()

	engine.Use(s.Logger)
}
