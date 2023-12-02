package user

import (
	"context"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

type Repository interface {
	CreateDialog(context.Context, *entities.Dialog) error
	EditDialog(context.Context, *entities.Dialog) error
	GetDialogByTalkerID(context.Context, int64) (*entities.Dialog, error)

	GetCategoriesByType(context.Context, int) ([]entities.Category, error)

	GetSaleProducts(context.Context) ([]entities.Product, error)
	GetProducersByCategory(context.Context, string, int) ([]string, error)
	GetProductsByProducer(context.Context, string, int) ([]entities.Product, error)
	GetProduct(context.Context, string, string, int) (*entities.Product, error)
	GetProductWithoutCategoryType(context.Context, string, string) (*entities.Product, error)

	GetProducersRepair(context.Context) ([]string, error)
	GetModelsRepair(context.Context, string) ([]string, error)
	GetRepairs(context.Context, string) ([]entities.Repair, error)
	GetRepairByModelAndName(context.Context, string, string) (*entities.Repair, error)

	GetUserIdsWhoManager(context.Context) ([]int64, error)

	CreateState(context.Context, *entities.State) error
	GetState(context.Context, string, int) (*entities.State, error)
}

type keyboardTextUserService struct {
	rep    Repository
	logger *zap.Logger
}

func ConfigureKeyboardTextUserService(app *app.App) {
	service := &keyboardTextUserService{
		rep:    app.GetRepository(),
		logger: app.GetLogger(),
	}
	engine := app.GetEngine()

	engine.Message("Главное меню", service.HomeMenu)
	engine.Message("Завершить диалог", service.ExitFromDialog)

	engine.Message("Акции", service.Sales)
	engine.Message("Новые устройства", service.Categories(0))
	engine.Message("Б/У устройства", service.Categories(1))
	engine.Message("Трейд-ин", service.TradeIn)
	engine.Message("Ремонт", service.CategoriesRepair)
	engine.Message("Связаться с менеджером", service.Manager)

	engine.MessageState("producers_0", service.Producers(0))
	engine.MessageState("producers_1", service.Producers(1))
	engine.MessageState("sales", service.SaleProduct)

	engine.GroupState("products_0", func(group bot.Router) {
		group.Message("Назад к категориям", service.Categories(0))
		group.MessageAny(service.Products(0))
	})

	engine.GroupState("products_1", func(group bot.Router) {
		group.Message("Назад к категориям", service.Categories(1))
		group.MessageAny(service.Products(1))
	})

	engine.GroupState("product_0", func(group bot.Router) {
		group.Message("Назад к категориям", service.Categories(0))
		group.MessageAny(service.Product(0))
	})

	engine.GroupState("product_1", func(group bot.Router) {
		group.Message("Назад к категориям", service.Categories(1))
		group.MessageAny(service.Product(1))
	})

	engine.GroupState("trade_in", func(group bot.Router) {
		group.Message("Отмена", service.HomeMenu)
		group.MessageAny(service.TradeInMessage)
	})

	engine.GroupState("repair_models", func(group bot.Router) {
		group.Message("Назад к списку производителей", service.CategoriesRepair)
		group.MessageAny(service.ModelsRepair)
	})

	engine.GroupState("repair", func(group bot.Router) {
		group.Message("Назад к списку производителей", service.CategoriesRepair)
		group.MessageAny(service.Repair)
	})

	engine.GroupState("repair_product", func(group bot.Router) {
		group.Message("Назад к списку производителей", service.CategoriesRepair)
		group.MessageAny(service.RepairProduct)
	})
}
