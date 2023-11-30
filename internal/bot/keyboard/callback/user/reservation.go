package user

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardCallbackUserService) Reservation(ctx *bot.Context) {
	id, err := ctx.GetCallbackData()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetCallbackData: %w", err))
		ctx.AbortWithCallback(true, "Не удалось получить техническую информацию.")
		return
	}

	product, err := service.rep.GetProductByID(ctx, uuid.MustParse(id.(string)))
	if err != nil {
		ctx.AddError(fmt.Errorf("rep.GetProductByID: %w", err))
		ctx.AbortWithCallback(true, "Не удалось получить информацию о товаре.")
		return
	}

	exists, err := service.rep.ExistsReservationByProductID(ctx, product.ID)
	if err != nil {
		ctx.AddError(fmt.Errorf("rep.ExistsReservationByProductID: %w", err))
		ctx.AbortWithCallback(true, "Не удалось получить техническую информацию.")
		return
	} else if exists {
		ctx.AbortWithCallback(true, "Данный товар уже забронирован, попробуйте позднее.")
		return
	}

	user, ok := ctx.Get("user")
	if !ok {
		ctx.AddError(fmt.Errorf("error get user by ctx.Get"))
		ctx.AbortWithCallback(true, "Данный товар уже зарезервирован, попробуйте позднее.")
		return
	}

	reservation := entities.Reservation{
		UserID:    user.(*entities.User).ID,
		ProductID: product.ID,
	}

	if err = service.rep.CreateReservation(ctx, &reservation); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateReservation: %w", err))
		ctx.AbortWithCallback(true, "Произошла ошибка при создании резервации.")
	}

	managers, err := service.rep.GetUserIdsWhoManager(ctx)
	if err != nil {
		ctx.AbortWithMessage("Произошла ошибка при получении списка менеджеров, для создания заявки")
		return
	}

	managerText := fmt.Sprintf("<i>Зарезервирован товар</i>.\n\nИмя и фамилия: <b>%s %s</b>\nТег: <b>%s</b>\n\nИнформация о товаре: \n%s", ctx.From().FirstName, ctx.From().LastName, ctx.From().UserName, product.StringWithoutDescription())
	managerKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Профиль пользователя", fmt.Sprintf("tg://user?id=%d", ctx.From().ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Начать диалог", fmt.Sprintf("manager_dialog_start_first:%d", ctx.From().ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтвердить резерв", fmt.Sprintf("manager_a_reserv:%d:%s", reservation.UserID, reservation.ProductID)), tgbotapi.NewInlineKeyboardButtonData("Отменить резерв", fmt.Sprintf("manager_c_reserv:%d:%s", reservation.UserID, reservation.ProductID)),
		),
	)

	var success bool
	for _, manager := range managers {
		if _, err = ctx.MessageWithKeyboardOtherChat(manager, managerText, managerKeyboard); err != nil {
			ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboardOtherChat: %w", err))
			continue
		}

		success = true
	}

	if !success {
		ctx.AbortWithMessage("В данный момент нет свободных менеджеров, попробуйте позже.")
		return
	}

	text := fmt.Sprintf("%s\n\n✅ Товар зарезервирован, ожидайте информации от менеджера.", product.String())
	if err = ctx.Edit(text); err != nil {
		ctx.AddError(fmt.Errorf("ctx.EditWithKeyboard: %w", err))
	}
}
