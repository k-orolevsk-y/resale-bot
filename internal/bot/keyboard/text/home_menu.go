package text

import (
	"fmt"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (s *service) HomeMenu(ctx *bot.Context) {
	text := "Главное меню"
	keyboard := constants.MainKeyboard()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}

	ctx.MustClearState()
	ctx.Abort()
}
