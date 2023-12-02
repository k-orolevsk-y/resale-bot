package manager

import (
	"context"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

type Repository interface {
	EditDialog(context.Context, *entities.Dialog) error
	GetDialogByTalkerID(context.Context, int64) (*entities.Dialog, error)

	GetProducts(context.Context) ([]entities.Product, error)
	GetCategories(context.Context) ([]entities.Category, error)
	GetAllRepairs(context.Context) ([]entities.Repair, error)
	GetReservations(context.Context) ([]entities.ReservationWithAdditionalData, error)
}

type keyboardCallbackManagerService struct {
	rep    Repository
	logger *zap.Logger
}

func ConfigureKeyboardCallbackManagerService(app *app.App) {
	service := keyboardCallbackManagerService{
		rep:    app.GetRepository(),
		logger: app.GetLogger(),
	}
	engine := app.GetEngine()

	engine.Group("manager_", func(group bot.Router) {
		group.Use(service.ManagerAccess)

		group.Callback("manager_dialog_start", service.ManagerDialogStart)
		group.Callback("manager_reserv_products", service.CatalogReservationProducts)
		group.Callback("manager_category", service.CatalogCategories)
		group.Callback("manager_products", service.CatalogProducts)
		group.Callback("manager_repairs", service.RepairsList)
	})
}
