package commands

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (s *service) Manager(ctx *bot.Context) {
	u, ok := ctx.Get("user")
	if !ok {
		ctx.AddError(fmt.Errorf("error get user by ctx.Get"))
		ctx.AbortWithMessage("Не удалось проверить права доступа.")
		return
	}
	user := u.(*entities.User)

	if !user.IsManager {
		s.logger.Info("user without manager right try use command /manager", zap.Any("user", user))
		ctx.AbortWithMessage("У вас нет доступа.")
		return
	}

	text := "Панель менеджера"
	keyboard := constants.ManagerKeyboard()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_panel")
	}
	ctx.Abort()
}
