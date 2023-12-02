package manager

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardTextManagerService) RepairEditState(ctx *bot.Context) {
	stateID := fmt.Sprintf("edit_repair_%d", ctx.From().ID)
	ctx.Set("edit_repair_state_id", stateID)

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

	product, err := service.rep.GetRepairByID(ctx, uuid.MustParse(m["repair_id"].(string)))
	if err != nil {
		ctx.AddError(fmt.Errorf("rep.GetRepairByID: %w", err))
		ctx.MustClearState()
		ctx.AbortWithMessage("Не удалось получить информацию о типе ремонта.")
		return
	}

	ctx.Set("edit_repair", product)
	ctx.Set("edit_repair_state", data)
}

func (service *keyboardTextManagerService) RepairEditManage(ctx *bot.Context) {
	state := ctx.MustGet("edit_repair_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "menu" {
		return
	}

	repair := ctx.MustGet("edit_repair").(*entities.Repair)

	var (
		text     string
		keyboard interface{} = nil
	)

	switch ctx.GetMessage().Text {
	case "Изменить производителя":
		producers, err := service.rep.GetProducersRepair(ctx)
		if err != nil {
			producers = []string{}
		}

		text = "Выберите или введите производителя устройства:"
		keyboard = constants.ManagerNewProductArrString(producers)

		data["action"] = "producer"
	case "Изменить модель":
		models, err := service.rep.GetModelsRepair(ctx, repair.ProducerName)
		if err != nil {
			models = []string{}
		}

		text = "Выберите или введите модель товара:"
		keyboard = constants.ManagerNewProductArrString(models)

		data["action"] = "model"
	case "Изменить название":
		text = "Введите название ремонта:"
		keyboard = constants.ManagerExit()

		data["action"] = "name"
	case "Изменить описание":
		text = "Введите описание ремонта:"
		keyboard = constants.ManagerEmpty()

		data["action"] = "description"
	case "Изменить цену":
		text = "Введите цену ремонта:"
		keyboard = constants.ManagerExit()

		data["action"] = "price"
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

func (service *keyboardTextManagerService) RepairEditProducer(ctx *bot.Context) {
	state := ctx.MustGet("edit_repair_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "producer" {
		return
	}

	repair := ctx.MustGet("edit_repair").(*entities.Repair)
	repair.ProducerName = ctx.GetMessage().Text

	if err := service.rep.EditRepair(ctx, repair); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditRepair: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить информацию о ремонте.")
		return
	} else {
		ctx.Set("edit_repair", repair)
	}

	service.RepairEdit(ctx)
}

func (service *keyboardTextManagerService) RepairEditModel(ctx *bot.Context) {
	state := ctx.MustGet("edit_repair_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "model" {
		return
	}

	repair := ctx.MustGet("edit_repair").(*entities.Repair)
	repair.ModelName = ctx.GetMessage().Text

	if err := service.rep.EditRepair(ctx, repair); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditRepair: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить информацию о ремонте.")
		return
	} else {
		ctx.Set("edit_repair", repair)
	}

	service.RepairEdit(ctx)
}

func (service *keyboardTextManagerService) RepairEditName(ctx *bot.Context) {
	state := ctx.MustGet("edit_repair_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "name" {
		return
	}

	repair := ctx.MustGet("edit_repair").(*entities.Repair)
	repair.Name = ctx.GetMessage().Text

	if err := service.rep.EditRepair(ctx, repair); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditRepair: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить информацию о ремонте.")
		return
	} else {
		ctx.Set("edit_repair", repair)
	}

	service.RepairEdit(ctx)
}

func (service *keyboardTextManagerService) RepairEditDescription(ctx *bot.Context) {
	state := ctx.MustGet("edit_repair_state").(*entities.State)
	data := state.Data.(map[string]interface{})

	if data["action"] != "description" {
		return
	}

	var description sql.NullString
	if ctx.GetMessage().Text == "Убрать" {
		description = sql.NullString{Valid: false, String: ""}
	} else {
		description = sql.NullString{Valid: true, String: ctx.GetMessage().Text}
	}

	repair := ctx.MustGet("edit_repair").(*entities.Repair)
	repair.Description = description

	if err := service.rep.EditRepair(ctx, repair); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditRepair: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить информацию о ремонте.")
		return
	} else {
		ctx.Set("edit_repair", repair)
	}

	service.RepairEdit(ctx)
}

func (service *keyboardTextManagerService) RepairEditPrice(ctx *bot.Context) {
	state := ctx.MustGet("edit_repair_state").(*entities.State)
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

	repair := ctx.MustGet("edit_repair").(*entities.Repair)
	repair.Price = price

	if err = service.rep.EditRepair(ctx, repair); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditRepair: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить информацию о ремонте.")
		return
	} else {
		ctx.Set("edit_repair", repair)
	}

	service.RepairEdit(ctx)
}

func (service *keyboardTextManagerService) DeleteRepair(ctx *bot.Context) {
	state := ctx.MustGet("edit_repair_state").(*entities.State)
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

		ctx.AbortWithMessage("Нажмите кнопку ещё раз в течении 10 секунд, если подтверждаете удаление <b>ремонта</b>.")
		return
	}

	repair := ctx.MustGet("edit_repair").(*entities.Repair)
	if err := service.rep.DeleteRepairByID(ctx, repair.ID); err != nil {
		ctx.AddError(fmt.Errorf("rep.DeleteRepairByID: %w", err))
		ctx.AbortWithMessage("Не удалось удалить ремонт.")
		return
	}

	service.Repairs(ctx)
}

func (service *keyboardTextManagerService) RepairEdit(ctx *bot.Context) {
	stateAny, ok := ctx.Get("edit_repair_state")
	if ok {
		state := stateAny.(*entities.State)
		data := state.Data.(map[string]interface{})

		data["action"] = "menu"
		state.Data = data

		ctx.Set("edit_repair_state", state)

		if err := service.rep.EditState(ctx, state); err != nil {
			ctx.AddError(fmt.Errorf("rep.EditState: %w", err))
		}
	}
	repair := ctx.MustGet("edit_repair").(*entities.Repair)

	text := repair.StringForBot("")
	keyboard := constants.ManagerRepairKeyboard()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_repair")
	}

	ctx.Abort()
}
