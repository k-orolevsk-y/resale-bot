package middlewares

import (
	"time"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (s *service) Logger(ctx *bot.Context) {
	start := time.Now()

	ctx.Next()

	query := ctx.Query()
	state := ctx.State()
	method := ctx.Method()
	fromID := ctx.From().ID
	updateID := ctx.GetUpdateID()
	duration := time.Since(start)

	s.logger.Info(
		"new update",
		zap.String("query", query),
		zap.String("state", state),
		zap.String("method", method),
		zap.Int64("from_id", fromID),
		zap.Int("update_id", updateID),
		zap.Duration("duration", duration),
	)

	if err := ctx.Error(); err != nil {
		s.logger.Info("error in update", zap.Int("update_id", updateID), zap.Error(err))
	}
}
