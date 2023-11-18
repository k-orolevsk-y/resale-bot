package bot

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/config"
)

type (
	Engine struct {
		botAPI *tgbotapi.BotAPI
		wg     sync.WaitGroup

		callbackStorage CallbackStorage

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

	return &Engine{
		botAPI:          botAPI,
		callbackStorage: newCallbackMemStorage(),

		recovery: defaultRecovery,
	}, nil
}

func (engine *Engine) SetCallbackStorage(storage CallbackStorage) {
	engine.callbackStorage = storage
}

func (engine *Engine) Use(middlewares ...HandlerFunc) *Engine {
	engine.handle("*", regexFroAllStrings, middlewares)
	engine.rebuildNoRouteHandlers(middlewares)

	return engine
}

func (engine *Engine) Command(template string, handlers ...HandlerFunc) *Engine {
	engine.handle("command", engine.mustCompileRegexByTemplate(template), handlers)
	return engine
}

func (engine *Engine) Message(template string, handlers ...HandlerFunc) *Engine {
	engine.handle("message", engine.mustCompileRegexByTemplate(template), handlers)
	return engine
}

func (engine *Engine) Callback(template string, handlers ...HandlerFunc) *Engine {
	engine.handle("callback", engine.mustCompileRegexByTemplate(template), handlers)
	return engine
}

func (engine *Engine) CommandRegex(regex *regexp.Regexp, handlers ...HandlerFunc) *Engine {
	engine.handle("command", regex, handlers)
	return engine
}

func (engine *Engine) MessageRegex(regex *regexp.Regexp, handlers ...HandlerFunc) *Engine {
	engine.handle("message", regex, handlers)
	return engine
}

func (engine *Engine) CallbackRegex(regex *regexp.Regexp, handlers ...HandlerFunc) *Engine {
	engine.handle("callback", regex, handlers)
	return engine
}

func (engine *Engine) Handle(method string, regex *regexp.Regexp, handlers ...HandlerFunc) *Engine {
	engine.handle(method, regex, handlers)
	return engine
}

func (engine *Engine) NoRoute(handlers ...HandlerFunc) *Engine {
	engine.handlersNoRoute = append(engine.handlersNoRoute, handlers...)
	return engine
}

func (engine *Engine) Recovery(handler RecoveryFunc) *Engine {
	engine.recovery = handler
	return engine
}

func (engine *Engine) Run() {
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

func (engine *Engine) handle(method string, regex *regexp.Regexp, handlers HandlersChain) {
	tree, ok := engine.trees.get(method)
	if !ok {
		tree = &methodTree{method: method, nodes: make([]*node, 0)}
		engine.trees = append(engine.trees, tree)
	}

	root, ok := tree.get(regex)
	if !ok {
		root = new(node)
		root.regex = regex
		tree.nodes = append(tree.nodes, root)
	}
	root.addHandlers(handlers)
}

func (engine *Engine) rebuildNoRouteHandlers(handlers HandlersChain) {
	engine.handlersNoRoute = append(engine.handlersNoRoute, handlers...)
}

func (engine *Engine) mustCompileRegexByTemplate(template string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^%s$", template))
}

func (engine *Engine) serveUpdate(ctx *Context) {
	defer engine.wg.Done()

	handlers := engine.trees.compileHandlersChain(ctx.Method(), ctx.Query())
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

		if !ctx.IsAborted() {
			ctx.Abort()
		}
	}
}
