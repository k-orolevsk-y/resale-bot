package text

import (
	"fmt"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (s *service) Manager(ctx *bot.Context) {
	var successMessageIds []string

	managerText := fmt.Sprintf("Поступила заявка на связь с менеджером.\n\nИмя и фамилия: <b>%s %s</b>\nТег: <b>%s</b>", ctx.From().FirstName, ctx.From().LastName, ctx.From().UserName)
	managerKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Профиль пользователя", fmt.Sprintf("tg://user?id=%d", ctx.From().ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Начать диалог", fmt.Sprintf("manager_dialog_start:%d", ctx.From().ID)),
		),
	)

	for _, manager := range constants.Managers {
		message, err := ctx.MessageWithKeyboardOtherChat(manager, managerText, managerKeyboard)
		if err != nil {
			continue
		}

		successMessageIds = append(successMessageIds, fmt.Sprint(message.MessageID))
	}

	if len(successMessageIds) < 1 {
		ctx.AbortWithMessage("В данный момент нет свободных менеджеров, попробуйте позже.")
		return
	}

	text := "В течение нескольких минут к Вам подключится менеджер и ответит на все интересующие Вас вопросы."
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена", fmt.Sprintf("cancel_manager:%s", strings.Join(successMessageIds, ","))),
		),
	)

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}
