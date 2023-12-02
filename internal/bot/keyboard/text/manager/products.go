package manager

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardTextManagerService) NewProduct(ctx *bot.Context) {
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

	stateID := fmt.Sprintf("new_product_%d", ctx.From().ID)
	state := entities.State{
		ID:   stateID,
		Type: 4,
		Data: map[string]interface{}{
			"step": "category",
		},
	}

	if err = service.rep.CreateState(ctx, &state); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateState: %w", err))
		ctx.AbortWithMessage("Ошибка при сохранении промежуточных значений.")
		return
	}

	text := "Выберите категорию нового товара:"
	keyboard := constants.ManagerNewProductCategory(categories)

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_new_product")
	}

	ctx.Abort()
}

func (service *keyboardTextManagerService) NewProductState(ctx *bot.Context) {
	stateID := fmt.Sprintf("new_product_%d", ctx.From().ID)
	ctx.Set("product_state_id", stateID)

	data, err := service.rep.GetState(ctx, stateID, 4)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			ctx.AddError(fmt.Errorf("rep.GetState: %w", err))
			ctx.AbortWithMessage("Не удалось получить техническую информацию.")
		}

		return
	}

	ctx.Set("product_state", data)
}

func (service *keyboardTextManagerService) NewProductCategory(ctx *bot.Context) {
	stateAny, ok := ctx.Get("product_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get product_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["step"] != "category" {
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

	producers, err := service.rep.GetProducersByCategoryID(ctx, category.ID)
	if err != nil {
		producers = []string{}
	}

	data["step"] = "producer"
	data["category_id"] = category.ID.String()
	state.Data = data

	if err = service.rep.EditState(ctx, state); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

	text := "Выберите или введите производителя товара:"
	keyboard := constants.ManagerNewProductArrString(producers)

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}

func (service *keyboardTextManagerService) NewProductProducer(ctx *bot.Context) {
	stateAny, ok := ctx.Get("product_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get product_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["step"] != "producer" {
		return
	}

	data["step"] = "model"
	data["producer"] = ctx.GetMessage().Text
	state.Data = data

	if err := service.rep.EditState(ctx, state); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

	categoryID := uuid.MustParse(data["category_id"].(string))

	models, err := service.rep.GetModelsByCategoryIDAndProducer(ctx, categoryID, ctx.GetMessage().Text)
	if err != nil {
		models = []string{}
	}

	text := "Выберите или введите модель товара:"
	keyboard := constants.ManagerNewProductArrString(models)

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}

func (service *keyboardTextManagerService) NewProductModel(ctx *bot.Context) {
	stateAny, ok := ctx.Get("product_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get product_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["step"] != "model" {
		return
	}

	data["step"] = "additional"
	data["model"] = ctx.GetMessage().Text
	state.Data = data

	if err := service.rep.EditState(ctx, state); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

	text := "Введите атрибуты товара (например: 512GB Gold):"
	keyboard := constants.ManagerExit()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}

func (service *keyboardTextManagerService) NewProductAdditional(ctx *bot.Context) {
	stateAny, ok := ctx.Get("product_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get product_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["step"] != "additional" {
		return
	}

	data["step"] = "description"
	data["additional"] = ctx.GetMessage().Text
	state.Data = data

	if err := service.rep.EditState(ctx, state); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

	text := "Введите описание товара:"
	keyboard := constants.ManagerSkip()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}

func (service *keyboardTextManagerService) NewProductDescription(ctx *bot.Context) {
	stateAny, ok := ctx.Get("product_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get product_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["step"] != "description" {
		return
	}

	description := ctx.GetMessage().Text
	if description == "Пропустить" {
		description = ""
	}

	data["step"] = "price"
	data["description"] = description
	state.Data = data

	if err := service.rep.EditState(ctx, state); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

	text := "Введите цену товара:"
	keyboard := constants.ManagerExit()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}

func (service *keyboardTextManagerService) NewProductPrice(ctx *bot.Context) {
	stateAny, ok := ctx.Get("product_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get product_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["step"] != "price" {
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

	data["step"] = "is_sale"
	data["price"] = price
	state.Data = data

	if err = service.rep.EditState(ctx, state); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

	text := "Товар по акции?"
	keyboard := constants.ManagerNewProductIsSale()

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}

func (service *keyboardTextManagerService) NewProductIsSale(ctx *bot.Context) {
	stateAny, ok := ctx.Get("product_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get product_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["step"] != "is_sale" {
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

	data["step"] = "image"
	data["is_sale"] = isSale
	state.Data = data

	if err := service.rep.EditState(ctx, state); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

	text := "Отправьте фото товара"
	keyboard := constants.ManagerSkip()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}

func (service *keyboardTextManagerService) NewProductImage(ctx *bot.Context) {
	stateAny, ok := ctx.Get("product_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get product_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["step"] != "image" {
		return
	}

	var photoID sql.NullString
	if ctx.GetMessage().Text == "Пропустить" {
		photoID.Valid = false
	} else {
		photo := ctx.GetMessage().Photo
		if photo == nil {
			ctx.AbortWithMessage("Нужно отправить фото товара или нажать кнопку \"Пропустить\".")
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

	stateID := ctx.MustGet("product_state_id").(string)
	if err := service.rep.DeleteState(ctx, stateID, 4); err != nil {
		ctx.AddError(fmt.Errorf("rep.DeleteState: %w", err))
	}

	categoryID := uuid.MustParse(data["category_id"].(string))
	product := entities.Product{
		CategoryID:  categoryID,
		Producer:    data["producer"].(string),
		Model:       data["model"].(string),
		Additional:  data["additional"].(string),
		Description: data["description"].(string),
		Photo:       photoID,
		Price:       data["price"].(float64),
		IsSale:      data["is_sale"].(bool),
	}

	if err := service.rep.CreateProduct(ctx, &product); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateProduct: %w", err))
		ctx.AbortWithMessage("Не удалось создать товар.")
		return
	}

	text := fmt.Sprintf("Товар #%s успешно создан!", strings.Split(product.ID.String(), "-")[0])
	keyboard := constants.ManagerCatalogProductsKeyboard()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_catalog_products")
	}

	ctx.Abort()
}

func (service *keyboardTextManagerService) CancelNewProduct(ctx *bot.Context) {
	if _, ok := ctx.Get("product_state"); ok {
		stateID := ctx.MustGet("product_state_id").(string)

		if err := service.rep.DeleteState(ctx, stateID, 4); err != nil {
			ctx.AddError(fmt.Errorf("rep.DeleteState: %w", err))
		}
	}

	service.CatalogProducts(ctx)
}

func (service *keyboardTextManagerService) ProductsList(ctx *bot.Context) {
	products, err := service.rep.GetProducts(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Товаров нет.")
		} else {
			ctx.AddError(fmt.Errorf("rep.GetProducts: %w", err))
			ctx.AbortWithMessage("Не удалось получить товары.")
		}
		return
	}

	botInfo, err := ctx.GetBot()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetBot: %w", err))
		ctx.AbortWithMessage("Ошибка при получении технической информации.")
		return
	}
	botURL := fmt.Sprintf("https://t.me/%s?start=", botInfo.UserName)

	countOnPage := 3
	text := "Товары:\n"

	for i, product := range products {
		text += fmt.Sprintf("\n%s\n", product.StringForBot(botURL))

		if i >= (countOnPage - 1) {
			break
		}
	}
	keyboard := constants.PaginationKeyboard("manager_products", len(products), 0, countOnPage)

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(err)
	}
	ctx.Abort()
}
