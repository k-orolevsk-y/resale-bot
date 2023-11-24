package app

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (a *App) noRoute(ctx *bot.Context) {
	if ctx.Method() == "command" || ctx.Method() == "message" {
		text := "Неизвестная команда, возможно устарела техническая информация, отправляю в главное меню."
		keyboard := tgbotapi.NewReplyKeyboard(
			constants.MainKeyboard...,
		)

		if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
			a.logger.Error("error send message (no route)", zap.Error(err))
		}
		ctx.Abort()

	} else if ctx.Method() == "callback" {
		ctx.AbortWithCallback(true, "Скорее всего устарела техническая информация, повторите действия.")
	}
}
