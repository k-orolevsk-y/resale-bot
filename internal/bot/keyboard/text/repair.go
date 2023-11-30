package text

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (s *service) CategoriesRepair(ctx *bot.Context) {
	var keyboard tgbotapi.ReplyKeyboardMarkup

	categories, err := s.rep.GetCategoriesRepair(ctx)
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
		keyboard = constants.CategoryRepairKeyboard(categories)
	}

	if err = ctx.MessageWithKeyboard("Список производителей устройств для ремонта", keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("repair_models")
	}

	ctx.Abort()
}

func (s *service) ModelsRepair(ctx *bot.Context) {
	var keyboard tgbotapi.ReplyKeyboardMarkup

	models, err := s.rep.GetModelsRepair(ctx, ctx.GetMessage().Text)
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

func (s *service) Repair(ctx *bot.Context) {
	var keyboard tgbotapi.ReplyKeyboardMarkup

	repairs, err := s.rep.GetRepairs(ctx, ctx.GetMessage().Text)
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
	}

	state := entities.State{
		ID:   fmt.Sprintf("repair_product_%d", ctx.From().ID),
		Type: 2,
		Data: repairs[0].ModelID,
	}

	if err = s.rep.CreateState(ctx, &state); err != nil {
		ctx.AbortWithMessage("Ошибка назначения промежуточных значений.")
		return
	}

	if err = ctx.MessageWithKeyboard("Список возможных ремонтов", keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("repair_product")
	}

	ctx.Abort()
}

func (s *service) RepairProduct(ctx *bot.Context) {
	stateID := fmt.Sprintf("repair_product_%d", ctx.From().ID)

	state, err := s.rep.GetState(ctx, stateID, 2)
	if err != nil {
		ctx.AbortWithMessage("Не удалось получить промежуточную информацию.")
		return
	}

	modelID := uuid.MustParse(state.Data.(string))
	repairName := strings.Split(ctx.GetMessage().Text, " - ")[0]

	repair, err := s.rep.GetRepairWithModelAndCategory(ctx, modelID, repairName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Данный ремонт пока что не принимается, попробуйте позже.")
			return
		} else {
			ctx.AddError(fmt.Errorf("rep.GetRepairWithModelAndCategory: %w", err))
			ctx.AbortWithMessage("Не удалось получить информацию.")
			return
		}
	}

	keyboard := constants.RepairKeyboard(&repair.Repair)
	if err = ctx.MessageWithKeyboard(repair.String(), keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
}
