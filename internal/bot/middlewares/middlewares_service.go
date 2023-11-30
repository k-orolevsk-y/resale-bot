package middlewares

import (
	"context"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

type Repository interface {
	CreateUser(context.Context, *entities.User) error
	EditUser(context.Context, *entities.User) error
	GetUserByTgID(context.Context, int64) (*entities.User, error)
}

type service struct {
	rep    Repository
	logger *zap.Logger
}

func ConfigureMiddlewaresService(app *app.App) {
	s := &service{
		rep:    app.GetRepository(),
		logger: app.GetLogger(),
	}
	engine := app.GetEngine()

	engine.Use(s.Logger)
	engine.Use(s.Auth)
	engine.Use(s.Ban)
}
