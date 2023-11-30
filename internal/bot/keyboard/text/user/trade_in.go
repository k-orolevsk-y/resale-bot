package user

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardTextUserService) TradeIn(ctx *bot.Context) {
	text := "♻️\nДля оценки устройства напишите пожалуйста ответы на вопросы:\n\n1. Модель устройства, объем памяти?\n2. Когда и Где покупали?\n3. В каком состоянии внешне (есть ли сколы, вмятины на корпусе? Если имеются, приложите фото)\n4. Имеется ли комплект (Коробка/Адаптер/Лайтнинг/Наушники/ Документы о покупке)\n5. Был ли в ремонтах? Все ли работает?\n6. Процент износа аккумулятора (можно посмотреть в настройках)"
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отмена"),
		),
	)

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}

	ctx.MustSetState("trade_in")
	ctx.Abort()
}

func (service *keyboardTextUserService) TradeInMessage(ctx *bot.Context) {
	managers, err := service.rep.GetUserIdsWhoManager(ctx)
	if err != nil {
		ctx.AbortWithMessage("Произошла ошибка при получении списка менеджеров, для создания заявки")
		return
	}

	dialog := entities.Dialog{
		UserID: ctx.From().ID,
	}

	if err = service.rep.CreateDialog(ctx, &dialog); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateDialog: %w", err))
		ctx.AbortWithMessage("Произошла ошибка при создании заявки.")
		return
	}

	managerText := fmt.Sprintf("Поступила заявка на <i>трейд-ин</i>.\n\nИмя и фамилия: <b>%service %service</b>\nТег: <b>%service</b>", ctx.From().FirstName, ctx.From().LastName, ctx.From().UserName)
	managerKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Профиль пользователя", fmt.Sprintf("tg://user?id=%d", ctx.From().ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Начать диалог", fmt.Sprintf("manager_dialog_start:%d", ctx.From().ID)),
		),
	)

	var success bool
	for _, manager := range managers {
		msg, err := ctx.MessageWithKeyboardOtherChat(manager, managerText, managerKeyboard)
		if err != nil {
			continue
		}

		cfg := tgbotapi.NewCopyMessage(manager, ctx.Chat().ID, ctx.GetMessage().MessageID)
		cfg.ReplyToMessageID = msg.MessageID

		if _, err = ctx.MessageByConfig(cfg); err != nil {
			service.logger.Error("error copy message for manager", zap.Error(err))
		}

		success = true
	}

	if !success {
		ctx.AbortWithMessage("В данный момент нет свободных менеджеров, попробуйте позже.")
		return
	}

	text := "В течение нескольких минут к Вам подключится менеджер и ответит на все интересующие Вас вопросы."
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_manager"),
		),
	)

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
	service.HomeMenu(ctx)
}
