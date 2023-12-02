package manager

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardTextManagerService) Repairs(ctx *bot.Context) {
	text := "Панель ремонтов"
	keyboard := constants.ManagerRepairsKeyboard()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_repairs")
	}

	ctx.Abort()
}

func (service *keyboardTextManagerService) NewRepair(ctx *bot.Context) {
	producers, err := service.rep.GetProducersRepair(ctx)
	if err != nil {
		producers = []string{}
	}

	stateID := fmt.Sprintf("new_repair_%d", ctx.From().ID)
	state := entities.State{
		ID:   stateID,
		Type: 4,
		Data: map[string]interface{}{
			"step": "producer",
		},
	}

	if err = service.rep.CreateState(ctx, &state); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateState: %w", err))
		ctx.AbortWithMessage("Ошибка при сохранении промежуточных значений.")
		return
	}

	text := "Выберите или введите производителя устройства:"
	keyboard := constants.ManagerNewProductArrString(producers)

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_new_repair")
	}

	ctx.Abort()
}

func (service *keyboardTextManagerService) NewRepairState(ctx *bot.Context) {
	stateID := fmt.Sprintf("new_repair_%d", ctx.From().ID)
	ctx.Set("repair_state_id", stateID)

	data, err := service.rep.GetState(ctx, stateID, 4)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			ctx.AddError(fmt.Errorf("rep.GetState: %w", err))
			ctx.AbortWithMessage("Не удалось получить техническую информацию.")
		}

		return
	}

	ctx.Set("repair_state", data)
}

func (service *keyboardTextManagerService) NewRepairProducer(ctx *bot.Context) {
	stateAny, ok := ctx.Get("repair_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get repair_state: %v (%T)", stateAny, stateAny))
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

	models, err := service.rep.GetModelsRepair(ctx, ctx.GetMessage().Text)
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

func (service *keyboardTextManagerService) NewRepairModel(ctx *bot.Context) {
	stateAny, ok := ctx.Get("repair_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get repair_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["step"] != "model" {
		return
	}

	data["step"] = "name"
	data["model"] = ctx.GetMessage().Text
	state.Data = data

	if err := service.rep.EditState(ctx, state); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

	text := "Введите название ремонта:"
	keyboard := constants.ManagerExit()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}

func (service *keyboardTextManagerService) NewRepairName(ctx *bot.Context) {
	stateAny, ok := ctx.Get("repair_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get repair_state: %v (%T)", stateAny, stateAny))
		ctx.AbortWithMessage("Произошла потеря промежуточной информации, попробуйте заново.")
		return
	}

	state := stateAny.(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["step"] != "name" {
		return
	}

	data["step"] = "description"
	data["name"] = ctx.GetMessage().Text
	state.Data = data

	if err := service.rep.EditState(ctx, state); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

	text := "Введите описание ремонта:"
	keyboard := constants.ManagerSkip()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}

func (service *keyboardTextManagerService) NewRepairDescription(ctx *bot.Context) {
	stateAny, ok := ctx.Get("repair_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get repair_state: %v (%T)", stateAny, stateAny))
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

	text := "Введите цену ремонта:"
	keyboard := constants.ManagerExit()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}

func (service *keyboardTextManagerService) NewRepairPrice(ctx *bot.Context) {
	stateAny, ok := ctx.Get("repair_state")
	if !ok {
		ctx.AddError(fmt.Errorf("error get repair_state: %v (%T)", stateAny, stateAny))
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

	stateID := ctx.MustGet("repair_state_id").(string)
	if err = service.rep.DeleteState(ctx, stateID, 4); err != nil {
		ctx.AddError(fmt.Errorf("rep.DeleteState: %w", err))
	}

	repair := entities.Repair{
		ProducerName: data["producer"].(string),
		ModelName:    data["model"].(string),
		Name:         data["name"].(string),
		Description:  sql.NullString{Valid: data["description"].(string) != "", String: data["description"].(string)},
		Price:        price,
	}

	if err = service.rep.CreateRepair(ctx, &repair); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateRepair: %w", err))
		ctx.AbortWithMessage("Не удалось создать новый тип ремонта.")
		return
	}

	text := fmt.Sprintf("Тип ремонта #%s успешно создан!", strings.Split(repair.ID.String(), "-")[0])
	keyboard := constants.ManagerRepairsKeyboard()

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_repairs")
	}

	ctx.Abort()
}

func (service *keyboardTextManagerService) CancelNewRepair(ctx *bot.Context) {
	if _, ok := ctx.Get("repair_state"); ok {
		stateID := ctx.MustGet("repair_state_id").(string)

		if err := service.rep.DeleteState(ctx, stateID, 4); err != nil {
			ctx.AddError(fmt.Errorf("rep.DeleteState: %w", err))
		}
	}

	service.Repairs(ctx)
}

func (service *keyboardTextManagerService) RepairsList(ctx *bot.Context) {
	repairs, err := service.rep.GetAllRepairs(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Ремонтов нет.")
		} else {
			ctx.AddError(fmt.Errorf("rep.GetAllRepairs: %w", err))
			ctx.AbortWithMessage("Не удалось получить ремонты.")
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
	text := "Ремонты:\n"

	for i, repair := range repairs {
		text += fmt.Sprintf("\n%s\n", repair.StringForBot(botURL))

		if i >= (countOnPage - 1) {
			break
		}
	}
	keyboard := constants.PaginationKeyboard("manager_repairs", len(repairs), 0, countOnPage)

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(err)
	}
	ctx.Abort()
}
