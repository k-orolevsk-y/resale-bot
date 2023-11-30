package text

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/tools"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (s *service) ExitFromDialog(ctx *bot.Context) {
	dialog, err := s.rep.GetDialogByTalkerID(ctx, ctx.From().ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithAnswer("Диалог уже завершен.")
		} else {
			ctx.AddError(fmt.Errorf("rep.GetDialogByTalkerID: %w", err))
			ctx.AbortWithAnswer("Не удалось получить данные диалога.")
		}
		return
	}

	dialog.EndedAt = tools.ProtoTime(time.Now())
	if err = s.rep.EditDialog(ctx, dialog); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditDialog: %w", err))
		ctx.AbortWithAnswer("Не удалось отменить диалог.")
		return
	}

	var (
		userText    string
		managerText string
	)

	keyboard := constants.MainKeyboard()

	if dialog.UserID == ctx.From().ID {
		userText = "Диалог завершен."
		managerText = "Пользователь завершил диалог."
	} else {
		userText = "Менеджер завершил диалог."
		managerText = "Диалог завершен."
	}

	if _, err = ctx.MessageWithKeyboardOtherChat(dialog.UserID, userText, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.AnswerWithKeyboard: %w", err))
	}

	if _, err = ctx.MessageWithKeyboardOtherChat(dialog.ManagerID, managerText, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboardOtherChat: %w", err))
	}

	ctx.MustClearOtherUserState(dialog.UserID)
	ctx.MustClearOtherUserState(dialog.ManagerID)

	ctx.Abort()
}
