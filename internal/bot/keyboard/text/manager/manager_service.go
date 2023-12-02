package manager

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

type Repository interface {
	FindUser(context.Context, interface{}) (*entities.User, error)
	EditUser(context.Context, *entities.User) error

	FindCategory(context.Context, string, int) (*entities.Category, error)
	EditCategory(context.Context, *entities.Category) error
	CreateCategory(context.Context, *entities.Category) error
	GetCategories(context.Context) ([]entities.Category, error)
	GetCategoryByID(context.Context, uuid.UUID) (*entities.Category, error)
	DeleteCategory(context.Context, uuid.UUID) error

	EditReservation(context.Context, *entities.Reservation) error
	GetReservations(context.Context) ([]entities.ReservationWithAdditionalData, error)
	GetReservationsByUserID(context.Context, int64) ([]entities.ReservationWithAdditionalData, error)
	GetReservationByID(context.Context, uuid.UUID) (*entities.ReservationWithAdditionalData, error)
	DeleteReservationsByCategoryID(context.Context, uuid.UUID) error
	DeleteReservationsByProductID(context.Context, uuid.UUID) error

	CreateProduct(context.Context, *entities.Product) error
	EditProduct(context.Context, *entities.Product) error
	GetProducts(context.Context) ([]entities.Product, error)
	GetProducersByCategoryID(context.Context, uuid.UUID) ([]string, error)
	GetModelsByCategoryIDAndProducer(context.Context, uuid.UUID, string) ([]string, error)
	GetProductByID(context.Context, uuid.UUID) (*entities.Product, error)
	DeleteProductByID(context.Context, uuid.UUID) error
	DeleteProductsByCategoryID(context.Context, uuid.UUID) error

	CreateRepair(context.Context, *entities.Repair) error
	EditRepair(context.Context, *entities.Repair) error
	GetModelsRepair(context.Context, string) ([]string, error)
	GetProducersRepair(context.Context) ([]string, error)
	GetAllRepairs(context.Context) ([]entities.Repair, error)
	GetRepairByID(context.Context, uuid.UUID) (*entities.Repair, error)
	DeleteRepairByID(context.Context, uuid.UUID) error

	CreateState(context.Context, *entities.State) error
	EditState(context.Context, *entities.State) error
	GetState(context.Context, string, int) (*entities.State, error)
	DeleteState(context.Context, string, int) error
}

type keyboardTextManagerService struct {
	rep    Repository
	logger *zap.Logger
}

func ConfigureKeyboardTextManagerService(app *app.App) {
	service := &keyboardTextManagerService{
		rep:    app.GetRepository(),
		logger: app.GetLogger(),
	}
	engine := app.GetEngine()

	engine.GroupState("manager_panel", func(group bot.Router) {
		group.Use(service.ManagerAccess)

		group.Message("Пользователи", service.Users)
		group.Message("Каталог", service.Catalog)
		group.Message("Ремонты", service.Repairs)
	})

	engine.GroupState("manager_panel_users", func(group bot.Router) {
		group.Use(service.ManagerAccess)

		group.Message("Отмена", service.Panel)
		group.MessageAny(service.GetUser)
	})

	engine.GroupState("manager_panel_user", func(group bot.Router) {
		group.Use(service.ManagerAccess)

		group.Message("Список бронирований", service.GetUserID, service.GetUserListReservation)
		group.Message("Сменить статус блокировки", service.GetUserID, service.ChangeUserStatusBanned)
		group.Message("Сменить статус прав менеджера", service.GetUserID, service.ChangeUserStatusManager)
		group.Message("Вернуться в панель менеджера", service.Panel)
	})

	engine.GroupState("manager_catalog", func(group bot.Router) {
		group.Use(service.ManagerAccess)

		group.Message("Забронированные товары", service.CatalogReservationProducts)
		group.Message("Категории", service.CatalogCategories)
		group.Message("Товары", service.CatalogProducts)
		group.Message("Вернуться в панель менеджера", service.Panel)
	})

	engine.GroupState("manager_catalog_categories", func(group bot.Router) {
		group.Use(service.ManagerAccess)

		group.Message("Создать новую категорию", service.NewCategory)
		group.Message("Список категорий", service.CategoriesList)
		group.Message("Вернуться в панель каталога", service.Catalog)
	})

	engine.GroupState("manager_catalog_products", func(group bot.Router) {
		group.Use(service.ManagerAccess)

		group.Message("Добавить новый товар", service.NewProduct)
		group.Message("Список товаров", service.ProductsList)
		group.Message("Вернуться в панель каталога", service.Catalog)
	})

	engine.GroupState("manager_repairs", func(group bot.Router) {
		group.Use(service.ManagerAccess)

		group.Message("Добавить новый ремонт", service.NewRepair)
		group.Message("Список ремонтов", service.RepairsList)
		group.Message("Вернуться в панель менеджера", service.Panel)
	})

	engine.GroupState("manager_new_category", func(group bot.Router) {
		group.Use(service.ManagerAccess)
		group.Use(service.NewCategoryState)

		group.Message("Отмена", service.CancelNewCategory)
		group.MessageAny(service.NewCategoryName, service.NewCategoryType)
	})

	engine.GroupState("manager_new_product", func(group bot.Router) {
		group.Use(service.ManagerAccess)
		group.Use(service.NewProductState)

		group.Message("Отмена", service.CancelNewProduct)
		group.MessageAny(
			service.NewProductCategory,
			service.NewProductProducer,
			service.NewProductModel,
			service.NewProductAdditional,
			service.NewProductDescription,
			service.NewProductPrice,
			service.NewProductIsSale,
			service.NewProductImage,
		)
	})

	engine.GroupState("manager_new_repair", func(group bot.Router) {
		group.Use(service.ManagerAccess)
		group.Use(service.NewRepairState)

		group.Message("Отмена", service.CancelNewRepair)
		group.MessageAny(
			service.NewRepairProducer,
			service.NewRepairModel,
			service.NewRepairName,
			service.NewRepairDescription,
			service.NewRepairPrice,
		)
	})

	engine.GroupState("manager_category", func(group bot.Router) {
		group.Use(service.ManagerAccess)
		group.Use(service.CategoriesEditState)

		group.Message("Вернуться в панель категорий", service.CatalogCategories)

		group.Message("Изменить название", service.CategoriesEditName)
		group.Message("Изменить тип", service.CategoriesEditType)
		group.Message("Удалить", service.DeleteCategory)

		group.MessageAny(service.CategoriesEditDBName)
	})

	engine.GroupState("manager_reservation", func(group bot.Router) {
		group.Use(service.ManagerAccess)
		group.Use(service.ReservationEditState)

		group.Message("Изменить статус", service.ReservationEditStatus)
		group.Message("Связаться с пользователем", service.CategoriesEditType)

		group.Message("Отменён", service.ReservationEditDBStatus(-1))
		group.Message("Рассматривается", service.ReservationEditDBStatus(0))
		group.Message("Выполнен", service.ReservationEditDBStatus(1))

		group.Message("Назад", service.Reservation)
		group.Message("Вернуться в панель каталога", service.Catalog)
	})

	engine.GroupState("manager_product", func(group bot.Router) {
		group.Use(service.ManagerAccess)
		group.Use(service.ProductsEditState)

		group.Message("Отмена", service.ProductsEdit)
		group.Message("Удалить", service.DeleteProduct)
		group.Message("Вернуться в панель товаров", service.CatalogProducts)

		group.MessageAny(
			service.ProductsEditManage,
			service.ProductsEditCategory,
			service.ProductsEditProducer,
			service.ProductsEditModel,
			service.ProductsEditAdditional,
			service.ProductsEditDescription,
			service.ProductsEditPrice,
			service.ProductsEditSalePrice,
			service.ProductsEditIsSale,
			service.ProductsEditImage,
		)
	})

	engine.GroupState("manager_repair", func(group bot.Router) {
		group.Use(service.ManagerAccess)
		group.Use(service.RepairEditState)

		group.Message("Отмена", service.RepairEdit)
		group.Message("Удалить", service.DeleteRepair)
		group.Message("Вернуться в панель ремонтов", service.Repairs)

		group.MessageAny(
			service.RepairEditManage,
			service.RepairEditProducer,
			service.RepairEditModel,
			service.RepairEditName,
			service.RepairEditDescription,
			service.RepairEditPrice,
		)
	})
}
