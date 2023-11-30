package bot

import (
	"errors"
	"fmt"
	"regexp"
	"sync"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/config"
)

type (
	Engine struct {
		Router

		botAPI *tgbotapi.BotAPI
		wg     sync.WaitGroup

		stateStorage    Storage
		callbackStorage Storage

		trees           methodTrees
		handlersNoRoute HandlersChain
		recovery        RecoveryFunc
	}

	HandlersChain []HandlerFunc
	HandlerFunc   func(ctx *Context)
	RecoveryFunc  func(ctx *Context, err interface{})
)

func New(logger *zap.Logger) (*Engine, error) {
	if err := tgbotapi.SetLogger(NewBotLogger(logger)); err != nil {
		return nil, err
	}

	botAPI, err := tgbotapi.NewBotAPI(config.Config.TelegramToken)
	if err != nil {
		return nil, err
	}

	if !config.Config.ProductionMode {
		botAPI.Debug = true
	}

	engine := &Engine{
		botAPI:   botAPI,
		recovery: defaultRecovery,
	}
	engine.Router = &RouterGroup{engine: engine}

	return engine, nil
}

func (engine *Engine) UseDefaultStorages() {
	engine.stateStorage = newMemStorage()
	engine.callbackStorage = newMemStorage()
}

func (engine *Engine) SetStateStorage(storage Storage) {
	engine.stateStorage = storage
}

func (engine *Engine) SetCallbackStorage(storage Storage) {
	engine.callbackStorage = storage
}

func (engine *Engine) SetCommands(commands ...tgbotapi.BotCommand) error {
	cfg := tgbotapi.NewSetMyCommands(commands...)

	_, err := engine.botAPI.Request(cfg)
	return err
}

func (engine *Engine) Use(middlewares ...HandlerFunc) Router {
	engine.Router.Use(middlewares...)
	engine.rebuildNoRouteHandlers(middlewares)

	return engine.Router
}

func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	engine.handlersNoRoute = append(engine.handlersNoRoute, handlers...)
}

func (engine *Engine) Recovery(handler RecoveryFunc) {
	engine.recovery = handler
}

func (engine *Engine) Run() {
	if engine.stateStorage == nil || engine.callbackStorage == nil {
		panic(ErrWithoutStorages)
	}

	go func() {
		cfg := tgbotapi.NewUpdate(0)
		cfg.Timeout = 30

		updates := engine.botAPI.GetUpdatesChan(cfg)
		for update := range updates {
			ctx := engine.allocateContext(update)
			ctx.reset()

			engine.wg.Add(1)
			go engine.serveUpdate(ctx)
		}
	}()
}

func (engine *Engine) StopReceivingUpdates() {
	engine.botAPI.StopReceivingUpdates()
	engine.wg.Wait()
}

func (engine *Engine) allocateContext(update tgbotapi.Update) *Context {
	return &Context{engine: engine, update: update}
}

func (engine *Engine) rebuildNoRouteHandlers(handlers HandlersChain) {
	engine.handlersNoRoute = append(engine.handlersNoRoute, handlers...)
}

func (engine *Engine) mustCompileRegexByTemplate(template string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^%s$", template))
}

func (engine *Engine) serveUpdate(ctx *Context) {
	defer engine.wg.Done()

	handlers := engine.trees.compileHandlersChain(ctx.Method(), ctx.Query(), ctx.State())
	if len(handlers) > 0 {
		ctx.handlers = handlers
	} else {
		ctx.handlers = engine.handlersNoRoute
	}

	defer engine.serveRecovery(ctx)
	ctx.Next()
}

func (engine *Engine) serveRecovery(ctx *Context) {
	if engine.recovery == nil {
		return
	}

	if err := recover(); err != nil {
		engine.recovery(ctx, err)

		if e, ok := err.(error); ok && e != nil {
			ctx.AddError(e)
		} else {
			ctx.AddError(errors.New(fmt.Sprint(e)))
		}

		if !ctx.IsAborted() {
			ctx.Abort()
		}
	}
}
