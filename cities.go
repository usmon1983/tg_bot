package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendCityKeyboard(bot *tgbotapi.BotAPI, chatID int64) {
    cities := []string{"Dushanbe", "Khujand", "Kurgan-Tyube", "Kulyab", "Khorog", "Hisor", "Gharm"}

    var rows [][]tgbotapi.KeyboardButton
    for _, city := range cities {
        row := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(city))
        rows = append(rows, row)
    }

    keyboard := tgbotapi.NewReplyKeyboard(rows...)
    msg := tgbotapi.NewMessage(chatID, "Выберите город для прогноза погоды:")
    msg.ReplyMarkup = keyboard

    bot.Send(msg)
}