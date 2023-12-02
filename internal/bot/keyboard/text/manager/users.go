package manager

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardTextManagerService) Users(ctx *bot.Context) {
	text := "Введите никнейм или ID пользователя"
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отмена"),
		),
	)
	ctx.MustSetState("manager_panel_users")

	_ = ctx.MessageWithKeyboard(text, keyboard)
	ctx.Abort()
}

func (service *keyboardTextManagerService) GetUser(ctx *bot.Context) {
	user, err := service.rep.FindUser(ctx, ctx.GetMessage().Text)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Такой пользователь не зарегистрирован в боте.")
			return
		} else {
			ctx.AddError(fmt.Errorf("rep.FindUser: %w", err))
			ctx.AbortWithMessage("Не удалось найти пользователя, попробуйте позже.")
			return
		}
	}

	state := entities.State{
		ID:   fmt.Sprintf("manager_panel_user_%d", ctx.From().ID),
		Type: 4,
		Data: user.ID,
	}

	if err = service.rep.CreateState(ctx, &state); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateState: %w", err))
		ctx.AbortWithMessage("Не удалось записать техническую информацию.")
		return
	}

	text := fmt.Sprintf("Информация о пользователе:\n\n%s", user.String())
	keyboard := constants.ManagerUserKeyboard()

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}

	ctx.MustSetState("manager_panel_user")
	ctx.Abort()
}

func (service *keyboardTextManagerService) GetUserID(ctx *bot.Context) {
	stateID := fmt.Sprintf("manager_panel_user_%d", ctx.From().ID)

	state, err := service.rep.GetState(ctx, stateID, 4)
	if err != nil {
		ctx.AddError(fmt.Errorf("rep.GetState: %w", err))
		ctx.AbortWithMessage("Не удалось получить ID пользователя над которым выполняются действия.")
		return
	}

	userID := int64(state.Data.(float64))
	ctx.Set("user_id", userID)
}

func (service *keyboardTextManagerService) GetUserListReservation(ctx *bot.Context) {
	userID := ctx.MustGet("user_id").(int64)

	reservations, err := service.rep.GetReservationsByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Данный пользователь не делал бронирование товаров.")
			return
		} else {
			ctx.AddError(fmt.Errorf("rep.GetReservationsByUserID: %w", err))
			ctx.AbortWithMessage("Не удалось получить бронирования товаров этого пользователя.")
			return
		}
	}

	var buf bytes.Buffer
	buf.WriteString("Категория | Товар | Дата | Статус\n")

	for _, reservation := range reservations {
		statusText := "Выполнен"
		if reservation.Completed == 0 {
			statusText = "Рассматривается"
		} else if reservation.Completed == -1 {
			statusText = "Отменен"
		}
		dateText := reservation.CreatedAt.Format("02.01.2006 15:04:05")

		text := fmt.Sprintf("\n%s | %s | %s | %s", reservation.CategoryName, reservation.ProductFullName, dateText, statusText)
		buf.WriteString(text)
	}

	fileName := fmt.Sprintf("reservations_user%d_%d.txt", userID, time.Now().Unix())
	file := tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: buf.Bytes(),
	}

	cfg := tgbotapi.NewDocument(ctx.Chat().ID, file)
	cfg.Caption = fmt.Sprintf("Бронирование товаров пользователя <b>%d</b>", userID)
	cfg.ParseMode = "HTML"

	if _, err = ctx.MessageByConfig(cfg); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageByConfig: %w", err))
	}
	ctx.Abort()
}

func (service *keyboardTextManagerService) ChangeUserStatusBanned(ctx *bot.Context) {
	userID := ctx.MustGet("user_id").(int64)

	if userID == ctx.From().ID {
		ctx.AbortWithMessage("Нельзя выдать блокировку самому себе!")
		return
	}

	user, err := service.rep.FindUser(ctx, userID)
	if err != nil {
		ctx.AddError(fmt.Errorf("rep.FindUser: %w", err))
		ctx.AbortWithMessage("Не удалось получить пользователя.")
		return
	}

	user.IsBanned = !user.IsBanned
	if err = service.rep.EditUser(ctx, user); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditUser: %w", err))
		ctx.AbortWithMessage("Не удалось сменить статус блокировки пользователя.")
		return
	}

	var (
		userText     string
		userKeyboard interface{}

		managerText string
	)

	if user.IsBanned {
		userText = "Администрация выдала вам блокировку аккаунта."
		userKeyboard = tgbotapi.NewRemoveKeyboard(true)

		managerText = fmt.Sprintf("Пользователь %s заблокирован.", user.Tag)
	} else {
		userText = "Администрация сняла вам блокировку аккаунта."
		userKeyboard = constants.MainKeyboard()

		managerText = fmt.Sprintf("Пользователь %s разблокирован.", user.Tag)
	}

	_, _ = ctx.MessageWithKeyboardOtherChat(user.ID, userText, userKeyboard)
	ctx.AbortWithMessage(managerText)
}

func (service *keyboardTextManagerService) ChangeUserStatusManager(ctx *bot.Context) {
	userID := ctx.MustGet("user_id").(int64)

	if userID == ctx.From().ID {
		ctx.AbortWithMessage("Нельзя изменить статус прав менеджера у самого себя!")
		return
	}

	user, err := service.rep.FindUser(ctx, userID)
	if err != nil {
		ctx.AddError(fmt.Errorf("rep.FindUser: %w", err))
		ctx.AbortWithMessage("Не удалось получить пользователя.")
		return
	}

	user.IsManager = !user.IsManager
	if err = service.rep.EditUser(ctx, user); err != nil {
		ctx.AddError(fmt.Errorf("rep.EditUser: %w", err))
		ctx.AbortWithMessage("Не удалось сменить статус прав менеджера пользователя.")
		return
	}

	var (
		userText    string
		managerText string
	)

	if user.IsManager {
		userText = "Администрация выдала вам права менеджера.\n\nПанель менеджера: /manager"
		managerText = fmt.Sprintf("Пользователю %s выданы права менеджера.", user.Tag)
	} else {
		ctx.MustClearOtherUserState(user.ID)

		userText = "Администрация сняла с вас права менеджера."
		managerText = fmt.Sprintf("С пользователя %s сняты права менеджера.", user.Tag)
	}

	_, _ = ctx.MessageWithKeyboardOtherChat(user.ID, userText, constants.MainKeyboard())
	ctx.AbortWithMessage(managerText)
}
