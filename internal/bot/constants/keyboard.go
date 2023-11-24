package constants

import "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	MainKeyboard = [][]tgbotapi.KeyboardButton{
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
	}
)
