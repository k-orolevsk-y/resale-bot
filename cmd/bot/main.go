package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/commands"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/config"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/keyboard/callback"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/keyboard/text"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/messages"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/middlewares"
	repository "github.com/k-orolevsk-y/resale-bot/internal/bot/repository/postgres"
	database "github.com/k-orolevsk-y/resale-bot/pkg/database/postgres"
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

	db, err := database.New()
	if err != nil {
		logger.Panic("error initialized postgres", zap.Error(err))
	}
	rep := repository.New(db)

	bot, err := app.New(logger, rep)
	if err != nil {
		logger.Panic("error initialized bot", zap.Error(err))
	}
	logger.Info("initialized bot")

	middlewares.ConfigureMiddlewaresService(bot)
	commands.ConfigureCommandsService(bot)
	text.ConfigureKeyboardTextService(bot)
	callback.ConfigureKeyboardCallbackService(bot)
	messages.ConfigureKeyboardMessagesService(bot)

	bot.Run()

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("shutting down gracefully", zap.Any("signal", <-quitSignal))
	bot.Stop()
	logger.Info("successfully shutdown bot gracefully")
}
