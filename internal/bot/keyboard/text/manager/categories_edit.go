package manager

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardTextManagerService) CategoriesEditState(ctx *bot.Context) {
	stateID := fmt.Sprintf("edit_category_%d", ctx.From().ID)
	ctx.Set("edit_category_state_id", stateID)

	data, err := service.rep.GetState(ctx, stateID, 4)
	if err != nil {
		ctx.AddError(fmt.Errorf("rep.GetState: %w", err))
		ctx.MustClearState()
		ctx.AbortWithMessage("Не удалось получить техническую информацию.")
		return
	}

	ctx.Set("edit_category_state", data)
}

func (service *keyboardTextManagerService) CategoriesEditName(ctx *bot.Context) {
	stateAny, ok := ctx.Get("edit_category_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get edit_category_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
	data := state.Data.(map[string]interface{})

	data["action"] = "name"
	state.Data = data

	if err := service.rep.EditState(ctx, state); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

	text := "Введите новое название"
	keyboard := constants.ManagerExit()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}

	ctx.Abort()
}

func (service *keyboardTextManagerService) CategoriesEditDBName(ctx *bot.Context) {
	stateAny, ok := ctx.Get("edit_category_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get edit_category_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "name" {
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

	categoryID := uuid.MustParse(data["category_id"].(string))

	category, err := service.rep.GetCategoryByID(ctx, categoryID)
	if err != nil {
		ctx.AddError(fmt.Errorf("rep.GetCategoryByID: %w", err))
		ctx.AbortWithMessage("Не удалось получить информацию о категории.")
		return
	}

	category.Name = categoryName
	if err = service.rep.EditCategory(ctx, category); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditCategory: %w", err))
		ctx.AbortWithMessage("Не удалось отредактировать категорию.")
		return
	}

	data["action"] = "menu"
	state.Data = data

	if err = service.rep.EditState(ctx, state); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
		_ = ctx.Message("Не удалось сохранить промежуточную информацию.")
	}

	service.category(ctx, category)
}

func (service *keyboardTextManagerService) CategoriesEditType(ctx *bot.Context) {
	stateAny, ok := ctx.Get("edit_category_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get edit_category_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
	data := state.Data.(map[string]interface{})
	categoryID := uuid.MustParse(data["category_id"].(string))

	category, err := service.rep.GetCategoryByID(ctx, categoryID)
	if err != nil {
		ctx.AddError(fmt.Errorf("rep.GetCategoryByID: %w", err))
		ctx.AbortWithMessage("Не удалось получить информацию о категории.")
		return
	}

	if category.Type == 0 {
		category.Type = 1
	} else {
		category.Type = 0
	}

	if err = service.rep.EditCategory(ctx, category); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditCategory: %w", err))
		ctx.AbortWithMessage("Не удалось отредактировать категорию.")
		return
	}

	service.category(ctx, category)
}

func (service *keyboardTextManagerService) DeleteCategory(ctx *bot.Context) {
	stateAny, ok := ctx.Get("edit_category_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get edit_category_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
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

		ctx.AbortWithMessage("Нажмите кнопку ещё раз в течении 10 секунд, если подтверждаете удаление <b>всех товаров и этой категории</b>.")
		return
	}

	categoryID := uuid.MustParse(data["category_id"].(string))

	if err := service.rep.DeleteReservationsByCategoryID(ctx, categoryID); err != nil {
		ctx.AddError(fmt.Errorf("rep.DeleteReservationsByCategoryID: %w", err))
		ctx.AbortWithMessage("Не удалось удалить бронирования.")
		return
	}

	if err := service.rep.DeleteCategory(ctx, categoryID); err != nil {
		ctx.AddError(fmt.Errorf("rep.DeleteCategory: %w", err))
		ctx.AbortWithMessage("Не удалось удалить категорию.")
		return
	}

	if err := service.rep.DeleteProductsByCategoryID(ctx, categoryID); err != nil {
		ctx.AddError(fmt.Errorf("rep.DeleteProductsByCategoryID: %w", err))
		ctx.AbortWithMessage("Не удалось удалить товары из этой категории.")
		return
	}

	service.CatalogCategories(ctx)
}

func (service *keyboardTextManagerService) category(ctx *bot.Context, category *entities.Category) {
	text := fmt.Sprintf("Информация о категории:\n\n%s", category.StringForManager())
	keyboard := constants.ManagerCategoryKeyboard()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}

	ctx.Abort()
}
