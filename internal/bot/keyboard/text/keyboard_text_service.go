package text

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

	GetCategoriesByType(context.Context, int) ([]entities.Category, error)

	GetSaleProducts(context.Context) ([]entities.Product, error)
	GetProducersByCategory(context.Context, string, int) ([]string, error)
	GetProductsByProducer(context.Context, string, int) ([]entities.Product, error)
	GetProduct(context.Context, string, string, int) (*entities.Product, error)
	GetProductWithoutCategoryType(context.Context, string, string) (*entities.Product, error)

	GetCategoriesRepair(context.Context) ([]entities.CategoryRepair, error)
	GetModelsRepair(context.Context, string) ([]entities.ModelRepair, error)
	GetRepairs(context.Context, string) ([]entities.Repair, error)
	GetRepairWithModelAndCategory(context.Context, uuid.UUID, string) (*entities.RepairWithModelAndCategory, error)

	GetUserIdsWhoManager(context.Context) ([]int64, error)

	CreateState(context.Context, *entities.State) error
	GetState(context.Context, string, int) (*entities.State, error)
}

type service struct {
	rep    Repository
	logger *zap.Logger
}

func ConfigureKeyboardTextService(app *app.App) {
	s := &service{
		rep:    app.GetRepository(),
		logger: app.GetLogger(),
	}
	engine := app.GetEngine()

	engine.Message("Главное меню", s.HomeMenu)
	engine.Message("Завершить диалог", s.ExitFromDialog)

	engine.Message("Акции", s.Sales)
	engine.Message("Новые устройства", s.Categories(0))
	engine.Message("Б/У устройства", s.Categories(1))
	engine.Message("Трейд-ин", s.TradeIn)
	engine.Message("Ремонт", s.CategoriesRepair)
	engine.Message("Связаться с менеджером", s.Manager)

	engine.MessageState("producers_0", s.Producers(0))
	engine.MessageState("producers_1", s.Producers(1))
	engine.MessageState("sales", s.SaleProduct)

	engine.GroupState("products_0", func(group bot.Router) {
		group.Message("Назад к категориям", s.Categories(0))
		group.MessageAny(s.Products(0))
	})

	engine.GroupState("products_1", func(group bot.Router) {
		group.Message("Назад к категориям", s.Categories(1))
		group.MessageAny(s.Products(1))
	})

	engine.GroupState("product_0", func(group bot.Router) {
		group.Message("Назад к категориям", s.Categories(0))
		group.MessageAny(s.Product(0))
	})

	engine.GroupState("product_1", func(group bot.Router) {
		group.Message("Назад к категориям", s.Categories(1))
		group.MessageAny(s.Product(1))
	})

	engine.GroupState("trade_in", func(group bot.Router) {
		group.Message("Отмена", s.HomeMenu)
		group.MessageAny(s.TradeInMessage)
	})

	engine.GroupState("repair_models", func(group bot.Router) {
		group.Message("Назад к списку производителей", s.CategoriesRepair)
		group.MessageAny(s.ModelsRepair)
	})

	engine.GroupState("repair", func(group bot.Router) {
		group.Message("Назад к списку производителей", s.CategoriesRepair)
		group.MessageAny(s.Repair)
	})

	engine.GroupState("repair_product", func(group bot.Router) {
		group.Message("Назад к списку производителей", s.CategoriesRepair)
		group.MessageAny(s.RepairProduct)
	})
}
