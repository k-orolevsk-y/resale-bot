package constants

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

func MainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Акции"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Б/У устройства"), tgbotapi.NewKeyboardButton("Новые устройства"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Трейд-ин"), tgbotapi.NewKeyboardButton("Ремонт"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Связаться с менеджером"),
		),
	)
}

func ManagerKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Пользователи"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Товары"), tgbotapi.NewKeyboardButton("Ремонт"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Главное меню"),
		),
	)
}

func SalesKeyboard(products []entities.Product) tgbotapi.ReplyKeyboardMarkup {
	var row int
	keyboard := tgbotapi.NewReplyKeyboard()

	for i, product := range products {
		if len(keyboard.Keyboard) <= row {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow())
		}

		name := fmt.Sprintf("%s - %s", product.Model, product.Additional)
		keyboard.Keyboard[row] = append(keyboard.Keyboard[row], tgbotapi.NewKeyboardButton(name))

		if (i+1)%2 == 0 {
			row++
		}
	}

	keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Главное меню")))
	return keyboard
}

func CategoryKeyboard(categories []entities.Category) tgbotapi.ReplyKeyboardMarkup {
	var row int
	keyboard := tgbotapi.NewReplyKeyboard()

	for i, category := range categories {
		if len(keyboard.Keyboard) <= row {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow())
		}

		keyboard.Keyboard[row] = append(keyboard.Keyboard[row], tgbotapi.NewKeyboardButton(category.Name))

		if (i+1)%2 == 0 {
			row++
		}
	}

	keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Главное меню")))
	return keyboard
}

func ProducersKeyboard(producers []string) tgbotapi.ReplyKeyboardMarkup {
	var row int
	keyboard := tgbotapi.NewReplyKeyboard()

	for i, producer := range producers {
		if len(keyboard.Keyboard) <= row {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow())
		}

		keyboard.Keyboard[row] = append(keyboard.Keyboard[row], tgbotapi.NewKeyboardButton(producer))

		if (i+1)%2 == 0 {
			row++
		}
	}

	keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Назад к категориям")))
	return keyboard
}

func ProductsKeyboard(products []entities.Product) tgbotapi.ReplyKeyboardMarkup {
	var row int
	keyboard := tgbotapi.NewReplyKeyboard()

	for i, product := range products {
		if len(keyboard.Keyboard) <= row {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow())
		}

		name := fmt.Sprintf("%s - %s", product.Model, product.Additional)
		keyboard.Keyboard[row] = append(keyboard.Keyboard[row], tgbotapi.NewKeyboardButton(name))

		if (i+1)%2 == 0 {
			row++
		}
	}

	keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Назад к категориям")))
	return keyboard
}

func ProductKeyboard(botName string, product *entities.Product) tgbotapi.InlineKeyboardMarkup {
	reservation := fmt.Sprintf("reservation:%s", product.ID)
	share := fmt.Sprintf(" - посмотри на товар в этом боте.\n\nhttps://t.me/%s?start=product_%s", botName, product.ID)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Забронировать", reservation),
			tgbotapi.NewInlineKeyboardButtonSwitch("Поделиться", share),
		),
	)
}

func CategoryRepairKeyboard(categories []entities.CategoryRepair) tgbotapi.ReplyKeyboardMarkup {
	var row int
	keyboard := tgbotapi.NewReplyKeyboard()

	for i, category := range categories {
		if len(keyboard.Keyboard) <= row {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow())
		}

		keyboard.Keyboard[row] = append(keyboard.Keyboard[row], tgbotapi.NewKeyboardButton(category.Name))

		if (i+1)%2 == 0 {
			row++
		}
	}

	keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Главное меню")))
	return keyboard
}

func ModelsRepairKeyboard(models []entities.ModelRepair) tgbotapi.ReplyKeyboardMarkup {
	var row int
	keyboard := tgbotapi.NewReplyKeyboard()

	for i, model := range models {
		if len(keyboard.Keyboard) <= row {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow())
		}

		keyboard.Keyboard[row] = append(keyboard.Keyboard[row], tgbotapi.NewKeyboardButton(model.Name))

		if (i+1)%2 == 0 {
			row++
		}
	}

	keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Назад к списку производителей")))
	return keyboard
}

func RepairsKeyboard(repairs []entities.Repair) tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard()
	for _, repair := range repairs {
		name := fmt.Sprintf("%s - %.2f ₽", repair.Name, repair.Price)
		keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(name)))
	}

	keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Назад к списку производителей")))
	return keyboard
}

func RepairKeyboard(repair *entities.Repair) tgbotapi.InlineKeyboardMarkup {
	repairData := fmt.Sprintf("repair:%s", repair.ID)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отправить заявку на ремонт", repairData),
		),
	)
}
