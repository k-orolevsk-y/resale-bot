package commands

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

type Repository interface {
	FindUser(context.Context, interface{}) (*entities.User, error)
	CreateState(context.Context, *entities.State) error

	GetProductByID(context.Context, uuid.UUID) (*entities.Product, error)
	GetCategoryByID(context.Context, uuid.UUID) (*entities.Category, error)
	GetRepairByID(ctx context.Context, id uuid.UUID) (*entities.Repair, error)
	GetReservationByID(context.Context, uuid.UUID) (*entities.ReservationWithAdditionalData, error)
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

	engine.Command(
		"start",
		s.Start, s.StartProduct, s.StartManagerAccess,
		s.StartManagerUser, s.StartManagerReservation, s.StartManagerCategory,
		s.StartManagerProduct, s.StartManagerRepair, s.StartUnknownCommand,
	)
	engine.Command("manager", s.Manager)
}
