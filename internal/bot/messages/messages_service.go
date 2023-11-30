package messages

import (
	"context"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

type Repository interface {
	GetDialogByTalkerID(context.Context, int64) (*entities.Dialog, error)
}

type service struct {
	rep    Repository
	logger *zap.Logger
}

func ConfigureKeyboardMessagesService(app *app.App) {
	s := &service{
		rep:    app.GetRepository(),
		logger: app.GetLogger(),
	}
	engine := app.GetEngine()

	engine.GroupState("manager_dialog", func(group bot.Router) {
		group.MessageAny(s.Dialog)
	})

}
