package messages

import (
	"fmt"

	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (s *service) Dialog(ctx *bot.Context) {
	dialog, err := s.rep.GetDialogByTalkerID(ctx, ctx.From().ID)
	if err != nil {
		if err = ctx.ClearState(); err != nil {
			ctx.AddError(fmt.Errorf("ctx.ClearState: %w", err))
		}

		ctx.AbortWithMessage("Диалог с менеджером был сломан, попробуйте ещё раз.")
		return
	}

	talkerID := dialog.UserID
	if ctx.From().ID == talkerID {
		talkerID = dialog.ManagerID
	}

	if _, err = ctx.CopyMessage(talkerID, ctx.Chat().ID, ctx.GetMessage().MessageID); err != nil {
		ctx.AddError(fmt.Errorf("ctx.CopyMessage: %w", err))
	}
	ctx.Abort()
}
