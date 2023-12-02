package manager

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardTextManagerService) ReservationEditState(ctx *bot.Context) {
	stateID := fmt.Sprintf("edit_reservation_%d", ctx.From().ID)

	data, err := service.rep.GetState(ctx, stateID, 4)
	if err != nil {
		ctx.AddError(fmt.Errorf("rep.GetState: %w", err))
		ctx.MustClearState()
		ctx.AbortWithMessage("Не удалось получить техническую информацию.")
		return
	}

	reservationID := uuid.MustParse(data.Data.(string))

	reservation, err := service.rep.GetReservationByID(ctx, reservationID)
	if err != nil {
		ctx.AddError(fmt.Errorf("rep.GetReservationByID: %w", err))
		ctx.MustClearState()
		ctx.AbortWithMessage("Не удалось получить техническую информацию.")
		return
	}

	ctx.Set("edit_reservation", reservation)
}

func (service *keyboardTextManagerService) ReservationEditStatus(ctx *bot.Context) {
	reservation := ctx.MustGet("edit_reservation").(*entities.ReservationWithAdditionalData)

	text := "Выберите новый статус:"
	keyboard := constants.ManagerReservationEditKeyboard(&reservation.Reservation)

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}

	ctx.Abort()
}

func (service *keyboardTextManagerService) ReservationEditDBStatus(status int) bot.HandlerFunc {
	return func(ctx *bot.Context) {
		reservation := ctx.MustGet("edit_reservation").(*entities.ReservationWithAdditionalData)

		reservation.Reservation.Completed = status
		if err := service.rep.EditReservation(ctx, &reservation.Reservation); err != nil {
			ctx.AddError(fmt.Errorf("rep.EditReservation: %w", err))
			ctx.AbortWithMessage("Не удалось отредактировать бронирование.")
			return
		}

		text := fmt.Sprintf("У брони #%s по товару <code>%s</code> изменен статус на ", strings.Split(reservation.ID.String(), "-")[0], reservation.ProductFullName)
		switch status {
		case -1:
			text += "<b>отменён</b>"
		case 0:
			text += "<b>рассматривается</b>"
		case 1:
			text += "<b>выполнен</b>"
		}

		if _, err := ctx.MessageOtherChat(reservation.UserID, text); err != nil {
			ctx.AddError(fmt.Errorf("ctx.MessageOtherChat: %w", err))
		}

		service.Reservation(ctx)
	}
}

func (service *keyboardTextManagerService) Reservation(ctx *bot.Context) {
	reservation := ctx.MustGet("edit_reservation").(*entities.ReservationWithAdditionalData)

	botInfo, err := ctx.GetBot()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetBot: %w", err))
		ctx.AbortWithCallback(true, "Ошибка при получении технической информации.")
		return
	}
	botURL := fmt.Sprintf("https://t.me/%s?start=", botInfo.UserName)

	text := reservation.StringForManager(botURL)
	keyboard := constants.ManagerReservationKeyboard()

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}

	ctx.Abort()
}
