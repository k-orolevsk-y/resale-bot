package middlewares

import (
	"fmt"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (s *service) Ban(ctx *bot.Context) {
	u, ok := ctx.Get("user")
	if !ok {
		ctx.AddError(fmt.Errorf("failed to get user to check for ban [1]"))
		return
	}

	user, ok := u.(*entities.User)
	if !ok {
		ctx.AddError(fmt.Errorf("failed to get user to check for ban [2]"))
		return
	}

	if user.IsBanned {
		if ctx.Method() == "callback" {
			if err := ctx.Callback(true, "Вам выдана блокировка администрацией, Вы не можете пользоваться ботом."); err != nil {
				ctx.AddError(fmt.Errorf("ctx.Callback: %w", err))
			}
		} else {
			if err := ctx.Message("Вам выдана блокировка администрацией, Вы не можете пользоваться ботом."); err != nil {
				ctx.AddError(fmt.Errorf("ctx.Message: %w", err))
			}
		}

		ctx.Abort()
		return
	}
}
