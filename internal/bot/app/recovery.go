package app

import (
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (a *App) Recovery(ctx *bot.Context, err interface{}) {
	updateID := ctx.GetUpdateID()
	a.logger.Error("recovery", zap.Int("update_id", updateID), zap.Any("error", err))

	ctx.AbortWithMessage("Не удалось обработать действие, попробуйте позднее.")
}
