package text

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (s *service) Manager(ctx *bot.Context) {
	managers, err := s.rep.GetUserIdsWhoManager(ctx)
	if err != nil {
		ctx.AbortWithMessage("Произошла ошибка при получении списка менеджеров")
		return
	}

	dialog := entities.Dialog{
		UserID: ctx.From().ID,
	}

	if err = s.rep.CreateDialog(ctx, &dialog); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateDialog: %w", err))
		ctx.AbortWithMessage("Произошла ошибка при создании диалога.")
		return
	}

	managerText := fmt.Sprintf("Поступила заявка на <i>связь с менеджером</i>.\n\nИмя и фамилия: <b>%s %s</b>\nТег: <b>%s</b>", ctx.From().FirstName, ctx.From().LastName, ctx.From().UserName)
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

	text := "Заявка отправлена, в скором времени Вам напишет менеджер."
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_manager"),
		),
	)

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}
