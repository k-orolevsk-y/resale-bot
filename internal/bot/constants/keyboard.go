package constants

import (
	"fmt"
	"math"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

// USER KEYBOARDS

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

func CategoryRepairKeyboard(categories []string) tgbotapi.ReplyKeyboardMarkup {
	var row int
	keyboard := tgbotapi.NewReplyKeyboard()

	for i, category := range categories {
		if len(keyboard.Keyboard) <= row {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow())
		}

		keyboard.Keyboard[row] = append(keyboard.Keyboard[row], tgbotapi.NewKeyboardButton(category))

		if (i+1)%2 == 0 {
			row++
		}
	}

	keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Главное меню")))
	return keyboard
}

func ModelsRepairKeyboard(models []string) tgbotapi.ReplyKeyboardMarkup {
	var row int
	keyboard := tgbotapi.NewReplyKeyboard()

	for i, model := range models {
		if len(keyboard.Keyboard) <= row {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow())
		}

		keyboard.Keyboard[row] = append(keyboard.Keyboard[row], tgbotapi.NewKeyboardButton(model))

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

// MANAGER KEYBOARDS

func ManagerKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Пользователи"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Каталог"), tgbotapi.NewKeyboardButton("Ремонты"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Главное меню"),
		),
	)
}

func ManagerUserKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Список бронирований"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Сменить статус блокировки"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Сменить статус прав менеджера"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вернуться в панель менеджера"),
		),
	)
}

func ManagerCatalogKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Забронированные товары"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Категории"), tgbotapi.NewKeyboardButton("Товары"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вернуться в панель менеджера"),
		),
	)
}

func ManagerCatalogCategoriesKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Создать новую категорию"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Список категорий"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вернуться в панель каталога"),
		),
	)
}

func ManagerCatalogProductsKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить новый товар"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Список товаров"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вернуться в панель каталога"),
		),
	)
}

func ManagerRepairsKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить новый ремонт"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Список ремонтов"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вернуться в панель менеджера"),
		),
	)
}

func ManagerNewCategory() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отмена"),
		),
	)
}

func ManagerNewCategoryType() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Новые"), tgbotapi.NewKeyboardButton("Б/У"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отмена"),
		),
	)
}

func ManagerNewProductCategory(categories []entities.Category) tgbotapi.ReplyKeyboardMarkup {
	var row int
	keyboard := tgbotapi.NewReplyKeyboard()

	for i, category := range categories {
		if len(keyboard.Keyboard) <= row {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow())
		}

		var name string
		if category.Type == 0 {
			name = fmt.Sprintf("%s [Новые]", category.Name)
		} else {
			name = fmt.Sprintf("%s [Б/У]", category.Name)
		}

		keyboard.Keyboard[row] = append(keyboard.Keyboard[row], tgbotapi.NewKeyboardButton(name))

		if (i+1)%2 == 0 {
			row++
		}
	}

	keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Отмена")))
	return keyboard
}

func ManagerNewProductArrString(arr []string) tgbotapi.ReplyKeyboardMarkup {
	if len(arr) < 1 {
		return tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Отмена"),
			),
		)
	}

	var row int
	keyboard := tgbotapi.NewReplyKeyboard()

	for i, elem := range arr {
		if len(keyboard.Keyboard) <= row {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow())
		}

		keyboard.Keyboard[row] = append(keyboard.Keyboard[row], tgbotapi.NewKeyboardButton(elem))

		if (i+1)%2 == 0 {
			row++
		}
	}

	keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Отмена")))
	return keyboard
}

func ManagerSkip() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Пропустить"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отмена"),
		),
	)
}

func ManagerEmpty() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Убрать"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отмена"),
		),
	)
}

func ManagerExit() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отмена"),
		),
	)
}

func ManagerNewProductIsSale() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Да"), tgbotapi.NewKeyboardButton("Нет"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отмена"),
		),
	)
}

func ManagerCategoryKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить название"), tgbotapi.NewKeyboardButton("Изменить тип"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Удалить"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вернуться в панель категорий"),
		),
	)
}

func ManagerReservationKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить статус"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вернуться в панель каталога"),
		),
	)
}

func ManagerProductKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить категорию"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить производителя"), tgbotapi.NewKeyboardButton("Изменить модель"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить атрибуты"), tgbotapi.NewKeyboardButton("Изменить описание"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить цену"), tgbotapi.NewKeyboardButton("Изменить скидку"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить статус акции"), tgbotapi.NewKeyboardButton("Изменить фото"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Удалить"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вернуться в панель товаров"),
		),
	)
}

func ManagerRepairKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить производителя"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить модель"), tgbotapi.NewKeyboardButton("Изменить название"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить описание"), tgbotapi.NewKeyboardButton("Изменить цену"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Удалить"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вернуться в панель ремонтов"),
		),
	)
}

func ManagerReservationEditKeyboard(reservation *entities.Reservation) tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard()

	if reservation.Completed != -1 {
		kb.Keyboard = append(kb.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Отменён")))
	}

	if reservation.Completed != 0 {
		kb.Keyboard = append(kb.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Рассматривается")))
	}

	if reservation.Completed != 1 {
		kb.Keyboard = append(kb.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Выполнен")))
	}

	kb.Keyboard = append(kb.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Назад")))
	return kb
}

func PaginationKeyboard(data string, countItems, offset, countOnPage int) tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(),
	)

	if offset > 0 {
		kbData := fmt.Sprintf("%s:%d", data, offset-countOnPage)
		kb.InlineKeyboard[0] = append(kb.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("<-", kbData))
	} else {
		kb.InlineKeyboard[0] = append(kb.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("|", "noData"))
	}

	currentPage := (offset / countOnPage) + 1
	maxPage := int(math.Ceil(float64(countItems) / float64(countOnPage)))
	mainText := fmt.Sprintf("%d / %d", currentPage, maxPage)

	kb.InlineKeyboard[0] = append(kb.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData(mainText, "noData"))

	if countItems > (offset + countOnPage) {
		kbData := fmt.Sprintf("%s:%d", data, offset+countOnPage)
		kb.InlineKeyboard[0] = append(kb.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("->", kbData))
	} else {
		kb.InlineKeyboard[0] = append(kb.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("|", "noData"))
	}

	return kb
}
