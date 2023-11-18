package bot

import (
	"strings"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/config"
)

type Logger struct {
	logger *zap.Logger
}

func NewBotLogger(logger *zap.Logger) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) Println(v ...interface{}) {
	if config.Config.ProductionMode {
		l.logger.Sugar().Info(v...)
	} else {
		l.logger.Sugar().Debug(v...)
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	format = strings.TrimSuffix(format, "\n")

	if config.Config.ProductionMode {
		l.logger.Sugar().Infof(format, v...)
	} else {
		l.logger.Sugar().Debugf(format, v...)
	}
}
