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

func (s *service) StartManagerAccess(ctx *bot.Context) {
	u, ok := ctx.Get("user")
	if !ok {
		ctx.AddError(fmt.Errorf("error get user by ctx.Get"))
		ctx.AbortWithMessage("Ошибка проверки команды, попробуйте ещё раз.")
		return
	}
	user := u.(*entities.User)

	if !user.IsManager {
		ctx.AbortWithMessage("Неизвестная команда, попробуйте ещё раз.")
	}
}

func (s *service) StartManagerUser(ctx *bot.Context) {
	args := ctx.GetMessage().CommandArguments()
	if !strings.HasPrefix(args, "mu-") {
		return
	}
	userID := strings.ReplaceAll(args, "mu-", "")

	user, err := s.rep.FindUser(ctx, userID)
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

	if err = s.rep.CreateState(ctx, &state); err != nil {
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

func (s *service) StartManagerCategory(ctx *bot.Context) {
	args := ctx.GetMessage().CommandArguments()
	if !strings.HasPrefix(args, "mc-") {
		return
	}
	categoryIDString := strings.ReplaceAll(args, "mc-", "")

	categoryID, err := uuid.Parse(categoryIDString)
	if err != nil {
		ctx.AbortWithMessage("ID категории невалиден.")
		return
	}

	category, err := s.rep.GetCategoryByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Категория не найдена.")
			return
		} else {
			ctx.AddError(fmt.Errorf("rep.GetCategoryByID: %w", err))
			ctx.AbortWithMessage("Не удалось получить информацию о категории.")
			return
		}
	}

	stateID := fmt.Sprintf("edit_category_%d", ctx.From().ID)
	state := entities.State{
		ID:   stateID,
		Type: 4,
		Data: map[string]interface{}{
			"action":      "menu",
			"category_id": categoryID,
		},
	}

	if err = s.rep.CreateState(ctx, &state); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

	text := fmt.Sprintf("Информация о категории:\n\n%s", category.StringForManager())
	keyboard := constants.ManagerCategoryKeyboard()

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_category")
	}

	ctx.Abort()
}

func (s *service) StartManagerProduct(ctx *bot.Context) {
	args := ctx.GetMessage().CommandArguments()
	if !strings.HasPrefix(args, "mp-") {
		return
	}
	productIDString := strings.ReplaceAll(args, "mp-", "")

	productID, err := uuid.Parse(productIDString)
	if err != nil {
		ctx.AbortWithMessage("ID товара невалиден.")
		return
	}

	product, err := s.rep.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Товар не найден.")
			return
		} else {
			ctx.AddError(fmt.Errorf("rep.GetProductByID: %w", err))
			ctx.AbortWithMessage("Не удалось получить информацию о товаре.")
			return
		}
	}

	stateID := fmt.Sprintf("edit_product_%d", ctx.From().ID)
	state := entities.State{
		ID:   stateID,
		Type: 4,
		Data: map[string]interface{}{
			"action":     "menu",
			"product_id": productID.String(),
		},
	}

	if err = s.rep.CreateState(ctx, &state); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

	botInfo, err := ctx.GetBot()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetBot: %w", err))
		ctx.AbortWithCallback(true, "Ошибка при получении технической информации.")
		return
	}
	botURL := fmt.Sprintf("https://t.me/%s?start=", botInfo.UserName)

	text := fmt.Sprintf("Информация о товаре:\n\n%s", product.StringForManager(botURL))
	keyboard := constants.ManagerProductKeyboard()

	if product.Photo.Valid {
		cfg := tgbotapi.NewPhoto(ctx.Chat().ID, tgbotapi.FileID(product.Photo.String))
		if _, err = ctx.MessageByConfig(cfg); err != nil {
			ctx.AddError(fmt.Errorf("ctx.MessageByConfig: %w", err))
		}
	}

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_product")
	}

	ctx.Abort()
}

func (s *service) StartManagerReservation(ctx *bot.Context) {
	args := ctx.GetMessage().CommandArguments()
	if !strings.HasPrefix(args, "mr-") {
		return
	}
	reservationIDString := strings.ReplaceAll(args, "mr-", "")

	reservationID, err := uuid.Parse(reservationIDString)
	if err != nil {
		ctx.AbortWithMessage("ID брони невалиден.")
		return
	}

	reservation, err := s.rep.GetReservationByID(ctx, reservationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Бронь не найдена.")
			return
		} else {
			ctx.AddError(fmt.Errorf("rep.GetProductByID: %w", err))
			ctx.AbortWithMessage("Не удалось получить информацию о брони.")
			return
		}
	}

	stateID := fmt.Sprintf("edit_reservation_%d", ctx.From().ID)
	state := entities.State{
		ID:   stateID,
		Type: 4,
		Data: reservationID.String(),
	}

	if err = s.rep.CreateState(ctx, &state); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

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
	} else {
		ctx.MustSetState("manager_reservation")
	}

	ctx.Abort()
}

func (s *service) StartManagerRepair(ctx *bot.Context) {
	args := ctx.GetMessage().CommandArguments()
	if !strings.HasPrefix(args, "mrp-") {
		return
	}
	repairIDString := strings.ReplaceAll(args, "mrp-", "")

	repairID, err := uuid.Parse(repairIDString)
	if err != nil {
		ctx.AbortWithMessage("ID типа ремонта невалиден.")
		return
	}

	repair, err := s.rep.GetRepairByID(ctx, repairID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Тип ремонта не найден.")
			return
		} else {
			ctx.AddError(fmt.Errorf("rep.GetProductByID: %w", err))
			ctx.AbortWithMessage("Не удалось получить информацию о типе ремонта.")
			return
		}
	}

	stateID := fmt.Sprintf("edit_repair_%d", ctx.From().ID)
	state := entities.State{
		ID:   stateID,
		Type: 4,
		Data: map[string]interface{}{
			"action":    "menu",
			"repair_id": repairID.String(),
		},
	}

	if err = s.rep.CreateState(ctx, &state); err != nil {
		ctx.AddError(fmt.Errorf("rep.CreateState: %w", err))
		ctx.AbortWithMessage("Не удалось сохранить промежуточную информацию.")
		return
	}

	text := repair.StringForBot("")
	keyboard := constants.ManagerRepairKeyboard()

	if err = ctx.MessageWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(fmt.Errorf("ctx.MessageWithKeyboard: %w", err))
	} else {
		ctx.MustSetState("manager_repair")
	}

	ctx.Abort()
}
