package text

import (
	"github.com/k-orolevsk-y/resale-bot/internal/bot/app"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/keyboard/text/user"
)

func ConfigureKeyboardTextService(app *app.App) {
	user.ConfigureKeyboardTextUserService(app)
}
