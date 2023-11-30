package user

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardTextUserService) Categories(cType int) bot.HandlerFunc {
	return func(ctx *bot.Context) {
		var keyboard tgbotapi.ReplyKeyboardMarkup

		categories, err := service.rep.GetCategoriesByType(ctx, cType)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				keyboard = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("В данном разделе нет доступных категорий"),
					),
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Главное меню"),
					),
				)
			} else {
				ctx.AddError(fmt.Errorf("rep.GetCategoriesByType: %w", err))
				ctx.AbortWithMessage("Произошла ошибка при получении товаров, попробуйте позже.")
				return
			}
		} else {
			keyboard = constants.CategoryKeyboard(categories)
		}

		if err = ctx.MessageWithKeyboard("Список категорий", keyboard); err != nil {
			ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
		} else {
			ctx.MustSetState(fmt.Sprintf("producers_%d", cType))
		}

		ctx.Abort()
	}
}

func (service *keyboardTextUserService) Producers(cType int) bot.HandlerFunc {
	return func(ctx *bot.Context) {
		var keyboard tgbotapi.ReplyKeyboardMarkup

		producers, err := service.rep.GetProducersByCategory(ctx, ctx.GetMessage().Text, cType)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				keyboard = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("В данной категории нет доступных товаров"),
					),
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Назад к категориям"),
					),
				)
			} else {
				ctx.AddError(fmt.Errorf("rep.GetProducersByCategory: %w", err))
				ctx.AbortWithMessage("Не удалось получить товары.")
				return
			}
		} else {
			keyboard = constants.ProducersKeyboard(producers)
		}

		if err = ctx.MessageWithKeyboard("Список товаров", keyboard); err != nil {
			ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
		} else {
			ctx.MustSetState(fmt.Sprintf("products_%d", cType))
		}

		ctx.Abort()
	}
}

func (service *keyboardTextUserService) Products(cType int) bot.HandlerFunc {
	return func(ctx *bot.Context) {
		var keyboard tgbotapi.ReplyKeyboardMarkup

		products, err := service.rep.GetProductsByProducer(ctx, ctx.GetMessage().Text, cType)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				keyboard = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("В данной категории нет доступных товаров"),
					),
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Назад к категориям"),
					),
				)
			} else {
				ctx.AddError(fmt.Errorf("rep.GetProductsByProducer: %w", err))
				ctx.AbortWithMessage("Не удалось получить товары.")
				return
			}
		} else {
			keyboard = constants.ProductsKeyboard(products)
		}

		if err = ctx.MessageWithKeyboard("Список моделей", keyboard); err != nil {
			ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
		} else {
			ctx.MustSetState(fmt.Sprintf("product_%d", cType))
		}

		ctx.Abort()
	}
}

func (service *keyboardTextUserService) Product(cType int) bot.HandlerFunc {
	return func(ctx *bot.Context) {
		name := strings.Split(ctx.GetMessage().Text, " - ")
		if len(name) < 2 {
			name = append(name, "")
		}

		product, err := service.rep.GetProduct(ctx, name[0], name[1], cType)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.AbortWithMessage("Этот товар уже продали.")
				return
			} else {
				ctx.AddError(fmt.Errorf("rep.GetProduct: %w", err))
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
}
