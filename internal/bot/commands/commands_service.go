package commands

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

type Repository interface {
	GetProductByID(context.Context, uuid.UUID) (*entities.Product, error)
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

	engine.Command("start", s.Start, s.StartProduct)
	engine.Command("manager", s.Manager)
}
