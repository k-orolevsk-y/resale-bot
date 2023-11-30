package user

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardCallbackUserService) Repair(ctx *bot.Context) {
	stateID := fmt.Sprintf("repair_product_%d", ctx.From().ID)

	state, err := service.rep.GetState(ctx, stateID, 2)
	if err != nil {
		ctx.AbortWithCallback(true, "Не удалось получить промежуточную информацию.")
		return
	}

	cbData, err := ctx.GetCallbackData()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetCallbackData: %w", err))
		ctx.AbortWithCallback(true, "Не удалось получить промежуточную информацию.")
		return
	}

	modelID := uuid.MustParse(state.Data.(string))
	repairID := uuid.MustParse(cbData.(string))

	repair, err := service.rep.GetRepairWithModelAndCategoryByID(ctx, modelID, repairID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithCallback(true, "Данный ремонт пока что не принимается, попробуйте позже.")
			return
		} else {
			ctx.AddError(fmt.Errorf("rep.GetRepairWithModelAndCategoryByID: %w", err))
			ctx.AbortWithCallback(true, "Не удалось получить информацию.")
			return
		}
	}

	ctx.Set("repair", repair)
}

func (service *keyboardCallbackUserService) RepairStartDialog(ctx *bot.Context) {
	r, ok := ctx.Get("repair")
	if !ok {
		ctx.AddError(fmt.Errorf("error get repair by ctx.Get"))
		ctx.AbortWithCallback(true, "Не удалось получить промежуточную информацию.")
		return
	}
	repair := r.(*entities.RepairWithModelAndCategory)

	managers, err := service.rep.GetUserIdsWhoManager(ctx)
	if err != nil {
		ctx.AbortWithMessage("Произошла ошибка при получении списка менеджеров, для создания заявки")
		return
	}

	dialog := entities.Dialog{
		UserID: ctx.From().ID,
	}

	if err = service.rep.CreateDialog(ctx, &dialog); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateDialog: %w", err))
		ctx.AbortWithMessage("Произошла ошибка при создании заявки.")
		return
	}

	managerText := fmt.Sprintf("Поступила заявка на <i>ремонт</i>.\n\nИмя и фамилия: <b>%service %service</b>\nТег: <b>%service</b>\n\nИнформация о ремонте: \n%service", ctx.From().FirstName, ctx.From().LastName, ctx.From().UserName, repair.StringWithoutDescription())
	managerKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Профиль пользователя", fmt.Sprintf("tg://user?id=%d", ctx.From().ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Начать диалог", fmt.Sprintf("manager_dialog_start:%d", ctx.From().ID)),
		),
	)

	var success bool
	for _, manager := range managers {
		if _, err = ctx.MessageWithKeyboardOtherChat(manager, managerText, managerKeyboard); err != nil {
			continue
		}

		success = true
	}

	if !success {
		ctx.AbortWithMessage("В данный момент нет свободных менеджеров, попробуйте позже.")
		return
	}

	text := "Ваша заявка отправлена, ожидайте связи с менеджером."
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_manager"),
		),
	)

	if err = ctx.EditWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.EditWithKeyboard: %w", err))
	}
}
