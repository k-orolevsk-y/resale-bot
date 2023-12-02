package manager

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardTextManagerService) Catalog(ctx *bot.Context) {
	text := "Панель управления каталогом"
	keyboard := constants.ManagerCatalogKeyboard()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_catalog")
	}

	ctx.Abort()
}

func (service *keyboardTextManagerService) CatalogReservationProducts(ctx *bot.Context) {
	reservations, err := service.rep.GetReservations(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("В данный момент нет забронированных товаров.")
		} else {
			ctx.AddError(fmt.Errorf("rep.GetReservations: %w", err))
			ctx.AbortWithMessage("Ошибка при получении забронированных товаров.")
		}
		return
	}

	botInfo, err := ctx.GetBot()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetBot: %w", err))
		ctx.AbortWithMessage("Ошибка при получении технической информации.")
		return
	}
	botURL := fmt.Sprintf("https://t.me/%s?start=", botInfo.UserName)

	countOnPage := 3
	text := "Забронированные товары:\n"

	for i, reservation := range reservations {
		text += fmt.Sprintf("\n%s\n", reservation.StringForBot(botURL))

		if i >= (countOnPage - 1) {
			break
		}
	}
	keyboard := constants.PaginationKeyboard("manager_reserv_products", len(reservations), 0, countOnPage)

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(err)
	}
	ctx.Abort()
}

func (service *keyboardTextManagerService) CatalogCategories(ctx *bot.Context) {
	text := "Панель управления категориями"
	keyboard := constants.ManagerCatalogCategoriesKeyboard()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_catalog_categories")
	}

	ctx.Abort()
}

func (service *keyboardTextManagerService) CatalogProducts(ctx *bot.Context) {
	text := "Панель управления товарами"
	keyboard := constants.ManagerCatalogProductsKeyboard()

	if err := ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_catalog_products")
	}

	ctx.Abort()
}
