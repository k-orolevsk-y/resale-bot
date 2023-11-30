package user

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardTextUserService) Sales(ctx *bot.Context) {
	var keyboard tgbotapi.ReplyKeyboardMarkup

	products, err := service.rep.GetSaleProducts(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			keyboard = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Нет товаров по акций"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Главное меню"),
				),
			)
		} else {
			ctx.AddError(fmt.Errorf("rep.GetSaleProducts: %w", err))
			ctx.AbortWithMessage("Произошла ошибка при получении товаров по акции, попробуйте позже.")
			return
		}
	} else {
		keyboard = constants.SalesKeyboard(products)
	}

	if err = ctx.MessageWithKeyboard("Товары по акции", keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("sales")
	}

	ctx.Abort()
}

func (service *keyboardTextUserService) SaleProduct(ctx *bot.Context) {
	name := strings.Split(ctx.GetMessage().Text, " - ")
	if len(name) < 2 {
		name = append(name, "")
	}

	product, err := service.rep.GetProductWithoutCategoryType(ctx, name[0], name[1])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Товар закончился.")
			return
		} else {
			ctx.AddError(fmt.Errorf("rep.GetProductWithoutCategoryType: %w", err))
			ctx.AbortWithMessage("Не удалось получить товар.")
			return
		}
	}

	botInfo, err := ctx.GetBot()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetBot: %w", err))
		ctx.AbortWithMessage("Не удалось получить технические данные.")
		return
	}

	if product.Photo.Valid {
		cfg := tgbotapi.NewPhoto(ctx.Chat().ID, tgbotapi.FileID(product.Photo.String))
		cfg.ParseMode = "HTML"

		if _, err = ctx.MessageByConfig(cfg); err != nil {
			ctx.AddError(fmt.Errorf("ctx.MessageByConfig: %w", err))
		}
	}

	keyboard := constants.ProductKeyboard(botInfo.UserName, product)
	if err = ctx.MessageWithKeyboard(product.String(), keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}

	ctx.Abort()
}
