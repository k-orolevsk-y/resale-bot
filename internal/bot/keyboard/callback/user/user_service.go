package user

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

type Repository interface {
	CreateDialog(context.Context, *entities.Dialog) error
	EditDialog(context.Context, *entities.Dialog) error
	GetDialogByTalkerID(context.Context, int64) (*entities.Dialog, error)

	GetUserIdsWhoManager(context.Context) ([]int64, error)
	GetRepairWithModelAndCategoryByID(context.Context, uuid.UUID, uuid.UUID) (*entities.RepairWithModelAndCategory, error)
	GetState(context.Context, string, int) (*entities.State, error)

	CreateReservation(context.Context, *entities.Reservation) error
	ExistsReservationByProductID(context.Context, uuid.UUID) (bool, error)

	GetProductByID(context.Context, uuid.UUID) (*entities.Product, error)
}

type keyboardCallbackUserService struct {
	rep    Repository
	logger *zap.Logger
}

func ConfigureKeyboardCallbackUserService(app *app.App) {
	service := keyboardCallbackUserService{
		rep:    app.GetRepository(),
		logger: app.GetLogger(),
	}

	engine := app.GetEngine()

	engine.Callback("repair", service.Repair, service.RepairStartDialog)
	engine.Callback("reservation", service.Reservation)
}
