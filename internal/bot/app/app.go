package app

import (
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/repository/postgres"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/repository/states"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

type App struct {
	engine *bot.Engine
	logger *zap.Logger
	rep    *postgres.Pg
}

func New(logger *zap.Logger, rep *postgres.Pg) (*App, error) {
	engine, err := bot.New(logger)
	if err != nil {
		return nil, err
	}

	if err = engine.SetCommands(constants.Commands...); err != nil {
		return nil, err
	}

	app := &App{engine: engine, logger: logger, rep: rep}

	engine.SetStateStorage(states.New(rep, 0))
	engine.SetCallbackStorage(states.New(rep, 1))

	engine.NoRoute(app.noRoute)
	engine.Recovery(app.Recovery)

	return app, nil
}

func (a *App) Run() {
	a.engine.Run()
}

func (a *App) Stop() {
	a.engine.StopReceivingUpdates()
}

func (a *App) GetEngine() *bot.Engine {
	return a.engine
}

func (a *App) GetLogger() *zap.Logger {
	return a.logger
}

func (a *App) GetRepository() *postgres.Pg {
	return a.rep
}
