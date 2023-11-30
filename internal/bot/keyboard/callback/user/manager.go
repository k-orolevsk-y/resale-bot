package user

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/tools"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardCallbackUserService) CancelManager(ctx *bot.Context) {
	dialog, err := service.rep.GetDialogByTalkerID(ctx, ctx.From().ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithCallback(true, "Диалог уже завершен.")
		} else {
			ctx.AddError(fmt.Errorf("rep.GetDialogByTalkerID: %w", err))
			ctx.AbortWithCallback(true, "Не удалось получить данные диалога.")
		}
		return
	} else if dialog.ManagerID != 0 {
		ctx.AbortWithCallback(true, "Отменить диалог уже нельзя, он уже начался.")
		return
	}

	dialog.EndedAt = tools.ProtoTime(time.Now())
	if err = service.rep.EditDialog(ctx, dialog); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditDialog: %w", err))
		ctx.AbortWithCallback(true, "Не удалось отменить диалог.")
		return
	}

	if err = ctx.Edit("Заявка отменена"); err != nil {
		ctx.AddError(fmt.Errorf("ctx.Edit: %w", err))
	}
	ctx.AbortWithCallback(true, "Диалог отменен.")
}
