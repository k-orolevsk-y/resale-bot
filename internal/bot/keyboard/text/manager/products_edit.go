package manager

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardTextManagerService) ProductsEditState(ctx *bot.Context) {
	stateID := fmt.Sprintf("edit_product_%d", ctx.From().ID)
	ctx.Set("edit_product_state_id", stateID)

	data, err := service.rep.GetState(ctx, stateID, 4)
	if err != nil {
		ctx.AddError(fmt.Errorf("rep.GetState: %w", err))
		ctx.MustClearState()
		ctx.AbortWithMessage("Не удалось получить техническую информацию.")
		return
	}

	m, ok := data.Data.(map[string]interface{})
	if !ok {
		ctx.AddError(fmt.Errorf("error convert state data to map[string]interface{}"))
		ctx.MustClearState()
		ctx.AbortWithMessage("Не удалось получить техническую информацию.")
		return
	}

	product, err := service.rep.GetProductByID(ctx, uuid.MustParse(m["product_id"].(string)))
	if err != nil {
		ctx.AddError(fmt.Errorf("rep.GetProductByID: %w", err))
		ctx.MustClearState()
		ctx.AbortWithMessage("Не удалось получить информацию о товаре.")
		return
	}

	ctx.Set("edit_product", product)
	ctx.Set("edit_product_state", data)
}

func (service *keyboardTextManagerService) ProductsEditManage(ctx *bot.Context) {
	state := ctx.MustGet("edit_product_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "menu" {
		return
	}

	product := ctx.MustGet("edit_product").(*entities.Product)

	var (
		text     string
		keyboard interface{} = nil
	)

	switch ctx.GetMessage().Text {
	case "Изменить категорию":
		categories, err := service.rep.GetCategories(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.AbortWithMessage("Нет категорий для товара.")
				return
			} else {
				ctx.AddError(fmt.Errorf("rep.GetCategories: %w", err))
				ctx.AbortWithMessage("Произошла ошибка при получении категорий для товара.")
				return
			}
		}

		text = "Выберите новую категорию для товара:"
		keyboard = constants.ManagerNewProductCategory(categories)

		data["action"] = "category"
	case "Изменить производителя":
		producers, err := service.rep.GetProducersByCategoryID(ctx, product.CategoryID)
		if err != nil {
			producers = []string{}
		}

		text = "Выберите или введите производителя товара:"
		keyboard = constants.ManagerNewProductArrString(producers)

		data["action"] = "producer"
	case "Изменить модель":
		models, err := service.rep.GetModelsByCategoryIDAndProducer(ctx, product.CategoryID, product.Producer)
		if err != nil {
			models = []string{}
		}

		text = "Выберите или введите модель товара:"
		keyboard = constants.ManagerNewProductArrString(models)

		data["action"] = "model"
	case "Изменить атрибуты":
		text = "Введите атрибуты товара (например: 512GB Gold):"
		keyboard = constants.ManagerExit()

		data["action"] = "additional"
	case "Изменить описание":
		text = "Введите описание товара:"
		keyboard = constants.ManagerEmpty()

		data["action"] = "description"
	case "Изменить цену":
		text = "Введите цену товара:"
		keyboard = constants.ManagerExit()

		data["action"] = "price"
	case "Изменить скидку":
		text = "Введите цену со скидкой для товара:"
		keyboard = constants.ManagerEmpty()

		data["action"] = "sale_price"
	case "Изменить статус акции":
		text = "Товар по акции?"
		keyboard = constants.ManagerNewProductIsSale()

		data["action"] = "is_sale"
	case "Изменить фото":
		text = "Отправьте фото товара"
		keyboard = constants.ManagerSkip()

		data["action"] = "image"
	}

	if text == "" {
		ctx.AbortWithMessage("Неверное действие")
		return
	}

	state.Data = data
	if err := service.rep.EditState(ctx, state); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

	if keyboard == nil {
		if err := ctx.Message(text); err != nil {
			ctx.AddError(fmt.Errorf("ctx.Message: %w", err))
		}
	} else {
		if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
			ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
		}
	}

	ctx.Abort()
}

func (service *keyboardTextManagerService) ProductsEditCategory(ctx *bot.Context) {
	state := ctx.MustGet("edit_product_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "category" {
		return
	}

	categoryType := -1
	messageText := ctx.GetMessage().Text

	if strings.Contains(messageText, "Новые") {
		categoryType = 0
	} else if strings.Contains(messageText, "Б/У") {
		categoryType = 1
	}
	categoryName := strings.Split(messageText, " [")[0]

	category, err := service.rep.FindCategory(ctx, categoryName, categoryType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Такой категории нет, выберите из клавиатуры.")
		} else {
			ctx.AddError(fmt.Errorf("rep.FindCategory: %w", err))
			ctx.AbortWithMessage("Ошибка при поиске категории.")
		}
		return
	}

	product := ctx.MustGet("edit_product").(*entities.Product)
	product.CategoryID = category.ID

	if err = service.rep.EditProduct(ctx, product); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditProduct: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить информацию о товаре.")
		return
	} else {
		ctx.Set("edit_product", product)
	}

	service.ProductsEdit(ctx)
}

func (service *keyboardTextManagerService) ProductsEditProducer(ctx *bot.Context) {
	state := ctx.MustGet("edit_product_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "producer" {
		return
	}

	product := ctx.MustGet("edit_product").(*entities.Product)
	product.Producer = ctx.GetMessage().Text

	if err := service.rep.EditProduct(ctx, product); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditProduct: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить информацию о товаре.")
		return
	} else {
		ctx.Set("edit_product", product)
	}

	service.ProductsEdit(ctx)
}

func (service *keyboardTextManagerService) ProductsEditModel(ctx *bot.Context) {
	state := ctx.MustGet("edit_product_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "model" {
		return
	}

	product := ctx.MustGet("edit_product").(*entities.Product)
	product.Model = ctx.GetMessage().Text

	if err := service.rep.EditProduct(ctx, product); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditProduct: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить информацию о товаре.")
		return
	} else {
		ctx.Set("edit_product", product)
	}

	service.ProductsEdit(ctx)
}

func (service *keyboardTextManagerService) ProductsEditAdditional(ctx *bot.Context) {
	state := ctx.MustGet("edit_product_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "additional" {
		return
	}

	product := ctx.MustGet("edit_product").(*entities.Product)
	product.Additional = ctx.GetMessage().Text

	if err := service.rep.EditProduct(ctx, product); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditProduct: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить информацию о товаре.")
		return
	} else {
		ctx.Set("edit_product", product)
	}

	service.ProductsEdit(ctx)
}

func (service *keyboardTextManagerService) ProductsEditDescription(ctx *bot.Context) {
	state := ctx.MustGet("edit_product_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "description" {
		return
	}

	description := ctx.GetMessage().Text
	if description == "Убрать" {
		description = ""
	}

	product := ctx.MustGet("edit_product").(*entities.Product)
	product.Description = description

	if err := service.rep.EditProduct(ctx, product); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditProduct: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить информацию о товаре.")
		return
	} else {
		ctx.Set("edit_product", product)
	}

	service.ProductsEdit(ctx)
}

func (service *keyboardTextManagerService) ProductsEditPrice(ctx *bot.Context) {
	state := ctx.MustGet("edit_product_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "price" {
		return
	}

	price, err := strconv.ParseFloat(ctx.GetMessage().Text, 64)
	if err != nil {
		ctx.AbortWithMessage("Введена неверная цена.")
		return
	} else if price < 1 {
		ctx.AbortWithMessage("Цена не может быть меньше 1.")
		return
	}

	product := ctx.MustGet("edit_product").(*entities.Product)
	product.Price = price

	if err = service.rep.EditProduct(ctx, product); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditProduct: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить информацию о товаре.")
		return
	} else {
		ctx.Set("edit_product", product)
	}

	service.ProductsEdit(ctx)
}

func (service *keyboardTextManagerService) ProductsEditSalePrice(ctx *bot.Context) {
	state := ctx.MustGet("edit_product_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "sale_price" {
		return
	}

	product := ctx.MustGet("edit_product").(*entities.Product)
	if ctx.GetMessage().Text == "Убрать" {
		if product.OldPrice == 0 {
			ctx.AbortWithMessage("У товара нет скидки.")
			return
		}

		product.Price = product.OldPrice
		product.OldPrice = 0
	} else {
		salePrice, err := strconv.ParseFloat(ctx.GetMessage().Text, 64)
		if err != nil {
			ctx.AbortWithMessage("Введена неверная цена со скидкой.")
			return
		} else if salePrice < 1 {
			ctx.AbortWithMessage("Цена со скидкой не может быть меньше 1.")
			return
		}

		product.OldPrice = product.Price
		product.Price = salePrice
	}

	if err := service.rep.EditProduct(ctx, product); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditProduct: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить информацию о товаре.")
		return
	} else {
		ctx.Set("edit_product", product)
	}

	service.ProductsEdit(ctx)
}

func (service *keyboardTextManagerService) ProductsEditIsSale(ctx *bot.Context) {
	state := ctx.MustGet("edit_product_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "is_sale" {
		return
	}

	var isSale bool
	switch ctx.GetMessage().Text {
	case "Да":
		isSale = true
	case "Нет":
		isSale = false
	default:
		ctx.AbortWithMessage("Выберите вариант на клавиатуре.")
		return
	}

	product := ctx.MustGet("edit_product").(*entities.Product)
	product.IsSale = isSale

	if err := service.rep.EditProduct(ctx, product); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditProduct: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить информацию о товаре.")
		return
	} else {
		ctx.Set("edit_product", product)
	}

	service.ProductsEdit(ctx)
}

func (service *keyboardTextManagerService) ProductsEditImage(ctx *bot.Context) {
	state := ctx.MustGet("edit_product_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "image" {
		return
	}

	var photoID sql.NullString
	if ctx.GetMessage().Text == "Пропустить" {
		photoID.Valid = false
	} else {
		photo := ctx.GetMessage().Photo
		if photo == nil {
			ctx.AbortWithMessage("Нужно отправить фото товара или нажать кнопку \"Убрать\".")
			return
		}

		var (
			maxPhotoID string
			width      int
			height     int
		)

		for _, ph := range photo {
			if ph.Width >= width && ph.Height >= height {
				maxPhotoID = ph.FileID
			}
		}

		photoID.Valid = true
		photoID.String = maxPhotoID
	}

	product := ctx.MustGet("edit_product").(*entities.Product)
	product.Photo = photoID

	if err := service.rep.EditProduct(ctx, product); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditProduct: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить информацию о товаре.")
		return
	} else {
		ctx.Set("edit_product", product)
	}

	service.ProductsEdit(ctx)
}

func (service *keyboardTextManagerService) DeleteProduct(ctx *bot.Context) {
	state := ctx.MustGet("edit_product_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	var flagConfirmTime bool
	confirmTimeUnix, ok := data["confirm_delete_time"]

	if !ok {
		flagConfirmTime = true
	} else {
		confirmTime := time.Unix(int64(confirmTimeUnix.(float64)), 0)
		if confirmTime.Unix()+10 <= time.Now().Unix() {
			flagConfirmTime = true
		}
	}

	if flagConfirmTime {
		data["confirm_delete_time"] = time.Now().Unix()
		state.Data = data

		if err := service.rep.EditState(ctx, state); err != nil {
			ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
			ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
			return
		}

		ctx.AbortWithMessage("Нажмите кнопку ещё раз в течении 10 секунд, если подтверждаете удаление <b>товара</b>.")
		return
	}

	product := ctx.MustGet("edit_product").(*entities.Product)
	if err := service.rep.DeleteProductByID(ctx, product.ID); err != nil {
		ctx.AddError(fmt.Errorf("rep.DeleteProductsByID: %w", err))
		ctx.AbortWithMessage("Не удалось удалить товар.")
		return
	}

	if err := service.rep.DeleteReservationsByProductID(ctx, product.ID); err != nil {
		ctx.AddError(fmt.Errorf("rep.DeleteReservationsByProductID: %w", err))
		_ = ctx.Message("Не удалось удалить бронирования.")
	}

	service.CatalogProducts(ctx)
}

func (service *keyboardTextManagerService) ProductsEdit(ctx *bot.Context) {
	stateAny, ok := ctx.Get("edit_product_state")
	if ok {
		state := stateAny.(*entities.State)
		data := state.Data.(map[string]interface{})

		data["action"] = "menu"
		state.Data = data

		ctx.Set("edit_product_state", state)

		if err := service.rep.EditState(ctx, state); err != nil {
			ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
		}
	}

	product := ctx.MustGet("edit_product").(*entities.Product)

	botInfo, err := ctx.GetBot()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetBot: %w", err))
		ctx.AbortWithCallback(true, "Ошибка при получении технической информации.")
		return
	}
	botURL := fmt.Sprintf("https://t.me/%s?start=", botInfo.UserName)

	text := fmt.Sprintf("Информация о товаре:\n\n%s", product.StringForManager(botURL))
	keyboard := constants.ManagerProductKeyboard()

	if product.Photo.Valid {
		cfg := tgbotapi.NewPhoto(ctx.Chat().ID, tgbotapi.FileID(product.Photo.String))
		if _, err = ctx.MessageByConfig(cfg); err != nil {
			ctx.AddError(fmt.Errorf("ctx.MessageByConfig: %w", err))
		}
	}

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}

	ctx.Abort()
}
