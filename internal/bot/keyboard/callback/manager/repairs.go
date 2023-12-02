package manager

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/constants"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (service *keyboardCallbackManagerService) RepairsList(ctx *bot.Context) {
	defer ctx.CallbackDone()

	offsetString, err := ctx.GetCallbackData()
	if err != nil {
		ctx.AddError(fmt.Errorf("ctx.GetCallbackData: %w", err))
		ctx.AbortWithCallback(true, "Не удалось получить техническую информацию.")
	}
	offset, _ := strconv.Atoi(offsetString.(string))

	repairs, err := service.rep.GetAllRepairs(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithMessage("Ремонтов нет.")
		} else {
			ctx.AddError(fmt.Errorf("rep.GetAllRepairs: %w", err))
			ctx.AbortWithMessage("Не удалось получить ремонты.")
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
	text := "Ремонты:\n"

	for i, repair := range repairs {
		if i < offset {
			continue
		}

		text += fmt.Sprintf("\n%s\n", repair.StringForBot(botURL))

		if (i - offset + 1) >= countOnPage {
			break
		}
	}
	keyboard := constants.PaginationKeyboard("manager_repairs", len(repairs), offset, countOnPage)

	if err = ctx.EditWithKeyboard(text, keyboard); err != nil {
		ctx.AddError(err)
	}
	ctx.Abort()
}
