package user

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardTextUserService) CategoriesRepair(ctx *bot.Context) {
	var keyboard tgbotapi.ReplyKeyboardMarkup

	producers, err := service.rep.GetProducersRepair(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			keyboard = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("В данный момент нет прайс-листов для ремонта"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Главное меню"),
				),
			)
		} else {
			ctx.AddError(fmt.Errorf("rep.GetCategoriesRepair: %w", err))
			ctx.AbortWithMessage("Произошла ошибка при получении списков, попробуйте позже.")
			return
		}
	} else {
		keyboard = constants.CategoryRepairKeyboard(producers)
	}

	if err = ctx.MessageWithKeyboard("Список производителей устройств для ремонта", keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("repair_models")
	}

	ctx.Abort()
}

func (service *keyboardTextUserService) ModelsRepair(ctx *bot.Context) {
	var keyboard tgbotapi.ReplyKeyboardMarkup

	models, err := service.rep.GetModelsRepair(ctx, ctx.GetMessage().Text)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			keyboard = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("В данный момент мы не ремонтируем данного производителя"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Назад к списку производителей"),
				),
			)
		} else {
			ctx.AddError(fmt.Errorf("rep.GetModelsRepair: %w", err))
			ctx.AbortWithMessage("Произошла ошибка при получении списков, попробуйте позже.")
			return
		}
	} else {
		keyboard = constants.ModelsRepairKeyboard(models)
	}

	if err = ctx.MessageWithKeyboard("Список устройств для ремонта", keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("repair")
	}

	ctx.Abort()
}

func (service *keyboardTextUserService) Repair(ctx *bot.Context) {
	var keyboard tgbotapi.ReplyKeyboardMarkup

	repairs, err := service.rep.GetRepairs(ctx, ctx.GetMessage().Text)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			keyboard = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("В данный момент мы не ремонтируем данное устройство"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Назад к списку производителей"),
				),
			)
		} else {
			ctx.AddError(fmt.Errorf("rep.GetRepairs: %w", err))
			ctx.AbortWithMessage("Произошла ошибка при получении списков, попробуйте позже.")
			return
		}
	} else {
		keyboard = constants.RepairsKeyboard(repairs)

		state := entities.State{
			ID:   fmt.Sprintf("repair_product_%d", ctx.From().ID),
			Type: 2,
			Data: repairs[0].ModelName,
		}

		if err = service.rep.CreateState(ctx, &state); err != nil {
			ctx.AbortWithMessage("Ошибка назначения промежуточных значений.")
			return
		}
	}

	if err = ctx.MessageWithKeyboard("Список возможных ремонтов", keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("repair_product")
	}

	ctx.Abort()
}

func (service *keyboardTextUserService) RepairProduct(ctx *bot.Context) {
	stateID := fmt.Sprintf("repair_product_%d", ctx.From().ID)

	state, err := service.rep.GetState(ctx, stateID, 2)
	if err != nil {
		ctx.AbortWithMessage("Не удалось получить промежуточную информацию.")
		return
	}

	modelName := state.Data.(string)
	repairName := strings.Split(ctx.GetMessage().Text, " - ")[0]

	repair, err := service.rep.GetRepairByModelAndName(ctx, modelName, repairName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Данный ремонт пока что не принимается, попробуйте позже.")
			return
		} else {
			ctx.AddError(fmt.Errorf("rep.GetRepairByModelAndName: %w", err))
			ctx.AbortWithMessage("Не удалось получить информацию.")
			return
		}
	}

	keyboard := constants.RepairKeyboard(repair)
	if err = ctx.MessageWithKeyboard(repair.String(), keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
}
