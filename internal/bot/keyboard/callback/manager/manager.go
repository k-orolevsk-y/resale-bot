package manager

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/tools"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardCallbackManagerService) ManagerAccess(ctx *bot.Context) {
	u, ok := ctx.Get("user")
	if !ok {
		ctx.AddError(fmt.Errorf("error get user by ctx.Get"))
		ctx.AbortWithCallback(true, "Не удалось проверить права доступа.")
		return
	}
	user := u.(*entities.User)

	if !user.IsManager {
		service.logger.Info("user without manager right try use callback buttons", zap.Any("user", user))
		ctx.AbortWithCallback(true, "У вас нет доступа.")
	}
}

func (service *keyboardCallbackManagerService) ManagerDialogStart(ctx *bot.Context) {
	data, err := ctx.GetCallbackData()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetCallbackData: %w", err))
		ctx.AbortWithCallback(false, "Не удалось получить данные.")

		return
	}

	if _, err = service.rep.GetDialogByTalkerID(ctx, ctx.From().ID); err == nil {
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

	dialog, err := service.rep.GetDialogByTalkerID(ctx, userChat.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithCallback(true, "Диалог уже завершен.")
		} else {
			ctx.AddError(fmt.Errorf("rep.GetDialogByTalkerID: %w", err))
			ctx.AbortWithCallback(true, "Не удалось получить данные диалога.")
		}
		return
	} else if dialog.EndedAt != nil {
		ctx.AbortWithCallback(true, "Диалог уже завершен.")
		return
	}

	dialog.ManagerID = ctx.From().ID
	if err = service.rep.EditDialog(ctx, dialog); err != nil {
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
