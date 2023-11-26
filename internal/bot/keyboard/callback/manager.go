package callback

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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

	if _, err = s.rep.GetDialogByTalkerID(ctx, ctx.From().ID); err == nil {
		ctx.AbortWithCallback(true, "У вас уже есть диалог с пользователем, закончите его перед тем как начать другой.")
		return
	}

	userChat, err := ctx.GetChat(tools.MustInt64(data.(string)))
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetChat: %w", err))
		ctx.AbortWithCallback(true, "Пользователь запретил отправлять сообщения боту, диалог невозможен.")
		return
	} else if userChat.ID == ctx.From().ID {
		ctx.AbortWithCallback(true, "Нельзя начать диалог с самим собой.")
		return
	}

	dialog, err := s.rep.GetDialogByTalkerID(ctx, userChat.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithCallback(true, "Диалог уже завершен.")
		} else {
			ctx.AddError(fmt.Errorf("rep.GetDialogByTalkerID: %w", err))
			ctx.AbortWithCallback(true, "Не удалось получить данные диалога.")
		}
		return
	}

	dialog.ManagerID = ctx.From().ID
	if err = s.rep.EditDialog(ctx, dialog); err != nil {
		ctx.AddError(fmt.Errorf("ctx.EditDialog: %w", err))
		ctx.AbortWithCallback(true, "Не удалось обозначить вас как менеджера диалога.")
		return
	}

	ctx.MustSetState("manager_dialog")
	ctx.MustSetOtherUserState(userChat.ID, "manager_dialog")

	cancelKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Завершить диалог"),
		),
	)

	if _, err = ctx.MessageWithKeyboardOtherChat(userChat.ID, "Менеджер подключился", cancelKeyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageOtherChat: %w", err))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Профиль пользователя", fmt.Sprintf("tg://user?id=%d", dialog.UserID)),
		),
	)

	if err = ctx.EditKeyboard(keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.EditKeyboard: %w", err))
	}

	if err = ctx.MessageWithKeyboard("Диалог начался", cancelKeyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	ctx.Abort()
}

func (s *service) CancelManager(ctx *bot.Context) {
	dialog, err := s.rep.GetDialogByTalkerID(ctx, ctx.From().ID)
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
	if err = s.rep.EditDialog(ctx, dialog); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditDialog: %w", err))
		ctx.AbortWithCallback(true, "Не удалось отменить диалог.")
		return
	}

	if err = ctx.Edit(ctx.GetMessage().Text); err != nil {
		ctx.AddError(fmt.Errorf("ctx.Edit: %w", err))
	}
	ctx.AbortWithCallback(true, "Диалог отменен.")
}
