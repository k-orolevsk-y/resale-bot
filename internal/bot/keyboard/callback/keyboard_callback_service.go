package callback

import (
	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/keyboard/callback/manager"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/keyboard/callback/user"
)

func ConfigureKeyboardCallbackService(app *app.App) {
	user.ConfigureKeyboardCallbackUserService(app)
	manager.ConfigureKeyboardCallbackManagerService(app)
}
