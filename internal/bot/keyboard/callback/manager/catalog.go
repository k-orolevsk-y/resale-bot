package manager

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardCallbackManagerService) CatalogReservationProducts(ctx *bot.Context) {
	defer ctx.CallbackDone()

	offsetString, err := ctx.GetCallbackData()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetCallbackData: %w", err))
		ctx.AbortWithCallback(true, "Не удалось получить техническую информацию.")
	}
	offset, _ := strconv.Atoi(offsetString.(string))

	reservations, err := service.rep.GetReservations(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithCallback(true, "В данный момент нет забронированных товаров.")
		} else {
			ctx.AddError(fmt.Errorf("rep.GetReservations: %w", err))
			ctx.AbortWithCallback(true, "Ошибка при получении забронированных товаров.")
		}
		return
	}

	botInfo, err := ctx.GetBot()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetBot: %w", err))
		ctx.AbortWithCallback(true, "Ошибка при получении технической информации.")
		return
	}
	botURL := fmt.Sprintf("https://t.me/%s?start=", botInfo.UserName)

	countOnPage := 3
	text := "Забронированные товары:\n"

	for i, reservation := range reservations {
		if i < offset {
			continue
		}

		text += fmt.Sprintf("\n%s\n", reservation.StringForBot(botURL))

		if (i - offset + 1) >= countOnPage {
			break
		}
	}
	keyboard := constants.PaginationKeyboard("manager_reserv_products", len(reservations), offset, countOnPage)

	if err = ctx.EditWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(err)
	}
	ctx.Abort()
}

func (service *keyboardCallbackManagerService) CatalogCategories(ctx *bot.Context) {
	defer ctx.CallbackDone()

	offsetString, err := ctx.GetCallbackData()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetCallbackData: %w", err))
		ctx.AbortWithCallback(true, "Не удалось получить техническую информацию.")
	}
	offset, _ := strconv.Atoi(offsetString.(string))

	categories, err := service.rep.GetCategories(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Категорий товаров нет.")
		} else {
			ctx.AddError(fmt.Errorf("rep.GetCategories: %w", err))
			ctx.AbortWithMessage("Не удалось получить категории товаров.")
		}
		return
	}

	botInfo, err := ctx.GetBot()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetBot: %w", err))
		ctx.AbortWithCallback(true, "Ошибка при получении технической информации.")
		return
	}
	botURL := fmt.Sprintf("https://t.me/%s?start=", botInfo.UserName)

	countOnPage := 5
	text := "Категории товаров:\n"

	for i, category := range categories {
		if i < offset {
			continue
		}

		text += fmt.Sprintf("\n%s\n", category.StringForBot(botURL))

		if (i - offset + 1) >= countOnPage {
			break
		}
	}
	keyboard := constants.PaginationKeyboard("manager_category", len(categories), offset, countOnPage)

	if err = ctx.EditWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(err)
	}
	ctx.Abort()
}

func (service *keyboardCallbackManagerService) CatalogProducts(ctx *bot.Context) {
	defer ctx.CallbackDone()

	offsetString, err := ctx.GetCallbackData()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetCallbackData: %w", err))
		ctx.AbortWithCallback(true, "Не удалось получить техническую информацию.")
	}
	offset, _ := strconv.Atoi(offsetString.(string))

	products, err := service.rep.GetProducts(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Товаров нет.")
		} else {
			ctx.AddError(fmt.Errorf("rep.GetCategories: %w", err))
			ctx.AbortWithMessage("Не удалось получить товары.")
		}
		return
	}

	botInfo, err := ctx.GetBot()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetBot: %w", err))
		ctx.AbortWithCallback(true, "Ошибка при получении технической информации.")
		return
	}
	botURL := fmt.Sprintf("https://t.me/%s?start=", botInfo.UserName)

	countOnPage := 3
	text := "Товары:\n"

	for i, product := range products {
		if i < offset {
			continue
		}

		text += fmt.Sprintf("\n%s\n", product.StringForBot(botURL))

		if (i - offset + 1) >= countOnPage {
			break
		}
	}
	keyboard := constants.PaginationKeyboard("manager_products", len(products), offset, countOnPage)

	if err = ctx.EditWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(err)
	}
	ctx.Abort()
}
