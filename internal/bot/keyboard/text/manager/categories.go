package manager

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardTextManagerService) NewCategory(ctx *bot.Context) {
	text := "Введите название новой категории"
	keyboard := constants.ManagerNewCategory()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_new_category")
	}

	ctx.Abort()
}

func (service *keyboardTextManagerService) NewCategoryState(ctx *bot.Context) {
	stateID := fmt.Sprintf("new_category_%d", ctx.From().ID)
	ctx.Set("category_state_id", stateID)

	data, err := service.rep.GetState(ctx, stateID, 4)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			ctx.AddError(fmt.Errorf("rep.GetState: %w", err))
			ctx.AbortWithMessage("Не удалось получить техническую информацию.")
		}

		return
	}

	ctx.Set("category_state", data)
}

func (service *keyboardTextManagerService) NewCategoryName(ctx *bot.Context) {
	if _, ok := ctx.Get("category_state"); ok {
		return
	}

	categoryName := ctx.GetMessage().Text
	if len(categoryName) <= 2 || len(categoryName) > 64 {
		ctx.AbortWithMessage("Название категории должно быть минимум 2 символа и максимум 64 символа.")
		return
	} else if strings.Contains(categoryName, "[") || strings.Contains(categoryName, "]") {
		ctx.AbortWithMessage("В названии категории не может быть <i>[</i> или <i>]</i>.")
		return
	}
	stateID := ctx.MustGet("category_state_id").(string)

	state := entities.State{
		ID:   stateID,
		Type: 4,
		Data: categoryName,
	}

	if err := service.rep.CreateState(ctx, &state); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию, попробуйте ещё раз.")
		return
	}

	text := "Выберите тип категории"
	keyboard := constants.ManagerNewCategoryType()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}

func (service *keyboardTextManagerService) NewCategoryType(ctx *bot.Context) {
	var categoryType int
	switch ctx.GetMessage().Text {
	case "Новые":
		categoryType = 0
	case "Б/У":
		categoryType = 1
	default:
		ctx.AbortWithMessage("Выберите категорию из доступных в клавиатуре.")
		return
	}

	stateAny, ok := ctx.Get("category_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get category_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
	categoryName := state.Data.(string)

	category := entities.Category{
		Name: categoryName,
		Type: categoryType,
	}

	if err := service.rep.CreateCategory(ctx, &category); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateCategory: %w", err))
		ctx.AbortWithMessage("Не удалось создать категорию, попробуйте ещё раз.")
		return
	} else {
		stateID := ctx.MustGet("category_state_id").(string)

		if err = service.rep.DeleteState(ctx, stateID, 4); err != nil {
			ctx.AddError(fmt.Errorf("rep.DeleteState: %w", err))
		}
		ctx.MustClearState()
	}

	text := fmt.Sprintf("Категория успешно создана (#%s), отправляю в панель категорий.", strings.Split(category.ID.String(), "-")[0])
	keyboard := constants.ManagerCatalogCategoriesKeyboard()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_catalog_categories")
	}

	ctx.Abort()
}

func (service *keyboardTextManagerService) CancelNewCategory(ctx *bot.Context) {
	if _, ok := ctx.Get("category_state"); ok {
		stateID := ctx.MustGet("category_state_id").(string)

		if err := service.rep.DeleteState(ctx, stateID, 4); err != nil {
			ctx.AddError(fmt.Errorf("rep.DeleteState: %w", err))
		}
	}

	service.CatalogCategories(ctx)
}

func (service *keyboardTextManagerService) CategoriesList(ctx *bot.Context) {
	categories, err := service.rep.GetCategories(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Категорий товаров нет.")
		} else {
			ctx.AddError(fmt.Errorf("rep.GetCategories: %w", err))
			ctx.AbortWithMessage("Не удалось получить категории товаров.")
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

	countOnPage := 5
	text := "Категории товаров:\n"

	for i, category := range categories {
		text += fmt.Sprintf("\n%s\n", category.StringForBot(botURL))

		if i >= (countOnPage - 1) {
			break
		}
	}
	keyboard := constants.PaginationKeyboard("manager_category", len(categories), 0, countOnPage)

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(err)
	}
	ctx.Abort()
}
