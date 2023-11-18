package log

import (
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/config"
)

func New() (*zap.Logger, error) {
	if config.Config.ProductionMode {
		return zap.NewProduction()
	}

	return zap.NewDevelopment()
}
