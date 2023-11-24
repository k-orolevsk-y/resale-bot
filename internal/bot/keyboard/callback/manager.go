package callback

import (
	"fmt"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/tools"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (s *service) ManagerAccess(ctx *bot.Context) {
	if !tools.In(ctx.From().ID, constants.Managers) {
		ctx.AbortWithCallback(true, "У вас нет доступа к этим командами.")
	}
}

func (s *service) ManagerDialogStart(ctx *bot.Context) {
	data, err := ctx.GetCallbackData()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetCallbackData: %w", err))
		ctx.AbortWithCallback(false, "Не удалось получить данные.")

		return
	}

	_, err = ctx.GetChat(tools.MustInt64(data.(string)))
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetChat: %w", err))
		ctx.AbortWithCallback(true, "Пользователь запретил отправлять сообщения боту, диалог невозможен.")
		return
	}

	panic("todo...")
	//ctx.MessageOtherChat(userChat.ID, "Менеджер подключился")
}

func (s *service) CancelManager(ctx *bot.Context) {

}
