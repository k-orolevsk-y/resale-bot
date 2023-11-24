package commands

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (s *service) Start(ctx *bot.Context) {
	text := "Успешная регистрация!"
	keyboard := tgbotapi.NewReplyKeyboard(
		constants.MainKeyboard...,
	)

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}
