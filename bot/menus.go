package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var menuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Нова дисципліна"),
		tgbotapi.NewKeyboardButton("Мої дисципліни"),
	),
)

func getDisciplineMenu(disciplines []string) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton
	var row []tgbotapi.KeyboardButton
	for i, discipline := range disciplines {
		if i%4 == 0 && i != 0 {
			rows = append(rows, tgbotapi.NewKeyboardButtonRow(row...))
			row = []tgbotapi.KeyboardButton{}
			println("append")
		}
		row = append(row, tgbotapi.NewKeyboardButton(discipline))
	}

	if len(row) != 0 {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(row...))
	}

	return tgbotapi.NewReplyKeyboard(rows...)
}

func getOperationMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Нове завдання"),
			tgbotapi.NewKeyboardButton("Виставити оцінку"),
		),
	)
}
