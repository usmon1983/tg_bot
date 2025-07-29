package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	log.Println("DATABASE_URL =", os.Getenv("DATABASE_URL"))
	err := InitDB()
	if err != nil {
		log.Fatal("Ошибка при инициализации базы:", err)
	}

	bot, err := tgbotapi.NewBotAPI(TelegramBotToken)
	if err != nil {
		log.Fatal("Ошибка при подключении к Telegram:", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//menu := tgbotapi.NewMessage(0, "Меню")
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("➕ Добавить расход"),
			tgbotapi.NewKeyboardButton("📊 Статистика"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("💰 Сегодняшние расходы"),
			tgbotapi.NewKeyboardButton("💱 Курс валют"),
		),
	)

	for update := range updates {
		go handleUpdate(bot, update, keyboard)
	}
}
