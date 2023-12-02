package callback

import "github.com/k-orolevsk-y/resale-bot/pkg/bot"

func noData(ctx *bot.Context) {
	ctx.CallbackDone()
	ctx.Abort()
}
