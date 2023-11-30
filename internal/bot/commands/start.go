package commands

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (s *service) Start(ctx *bot.Context) {
	var text string

	user, _ := ctx.Get("user")
	args := ctx.GetMessage().CommandArguments()

	if user.(*entities.User).RegisteredAt.IsZero() {
		text = "Успешная регистрация!"

		if args == "" {
			defer ctx.Abort()
		}
	} else {
		if args != "" {
			return
		}

		text = "Главное меню"
	}

	keyboard := constants.MainKeyboard()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	}
}

func (s *service) StartProduct(ctx *bot.Context) {
	args := ctx.GetMessage().CommandArguments()
	if !strings.HasPrefix(args, "product_") {
		return
	}
	productIDString := strings.ReplaceAll(args, "product_", "")

	productID, err := uuid.Parse(productIDString)
	if err != nil {
		ctx.AbortWithMessage("Данная ссылка на товар невалидна.")
		return
	}

	product, err := s.rep.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Этот товар уже продали.")
			return
		} else {
			ctx.AddError(fmt.Errorf("rep.GetProductByID: %w", err))
			ctx.AbortWithMessage("Не удалось получить информацию о товаре.")
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
