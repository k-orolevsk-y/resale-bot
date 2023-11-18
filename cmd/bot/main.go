package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/config"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
	"github.com/k-orolevsk-y/resale-bot/pkg/log"
)

func main() {
	if err := config.ParseConfig(); err != nil {
		panic(err)
	}

	logger, err := log.New()
	if err != nil {
		panic(err)
	}

	logger.Debug("initialized logger")
	logger.Debug("parsed config", zap.Any("config", config.Config))

	b, err := bot.New(logger)
	if err != nil {
		logger.Panic("error start bot", zap.Error(err))
	}

	b.Run()
	logger.Info("success started bot")

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("shutting down gracefully", zap.Any("signal", <-quitSignal))
	b.StopReceivingUpdates()
	logger.Info("successfully shutdown bot gracefully")
}
