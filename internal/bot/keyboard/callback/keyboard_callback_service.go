package callback

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
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

	GetProductByID(ctx context.Context, productID uuid.UUID) (*entities.Product, error)
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

	engine.Callback("repair", s.Repair, s.RepairStartDialog)
	engine.Callback("reservation", s.Reservation)
}
