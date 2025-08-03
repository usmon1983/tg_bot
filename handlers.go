package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, keyboard tgbotapi.ReplyKeyboardMarkup) {
	if update.Message == nil {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	text := update.Message.Text
	userID := update.Message.From.ID
	currentDate := time.Now().Format("02/01/06")

	switch text {
	case "/start":
		msg.Text = "Добро пожаловать! Используй кнопки ниже 👇"
		msg.ReplyMarkup = keyboard

	case "➕ Добавить расход":
		msg.Text = "Введите расход в формате: сумма категория\nНапример: `250 еда`"
		msg.ParseMode = "Markdown"

	case "💰 Сегодняшние расходы":
		total, err := GetTodayTotal(userID)
		if err != nil {
			msg.Text = "Ошибка при подсчёте"
		} else {
			msg.Text = fmt.Sprintf("Сегодня вы потратили: %.2f сомонӣ.", total)
		}

	case "📊 Статистика":
		stats, err := GetWeeklyStats(userID)
		if err != nil || len(stats) == 0 {
			msg.Text = "Не удалось получить статистику."
			break
		}
		filePath := "stats.png"
		if err := CreateStatsChart(stats, filePath); err != nil {
			msg.Text = "Ошибка построения графика."
			break
		}
		photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(filePath))
		photo.Caption = "📊 Ваши расходы за 7 дней"
		bot.Send(photo)
		return

	case "💱 Курс валют":
		rates, err := GetExchangeRatesFromXML()
		ratesMap := make(map[string]float64)
		pairs := strings.Split(rates, ",")
		for _, pair := range pairs {
			// Удаляем пробелы и разделяем по ":"
			parts := strings.Split(strings.TrimSpace(pair), ":")
			if len(parts) == 2 {
				code := strings.TrimSpace(parts[0])
				valueStr := strings.TrimSpace(parts[1])
				value, err := strconv.ParseFloat(valueStr, 64)
				if err != nil {
					fmt.Println("Ошибка парсинга числа:", err)
					continue
				}
				ratesMap[code] = value
			}
		}
		if err != nil {
			msg.Text = "Ошибка при получении курса валют."
		} else {
			msg.Text = fmt.Sprintf("Курс НБТ на дату %s", currentDate)
			msg.Text = fmt.Sprintf(
				"Курс НБТ на дату %s\n"+
					"🇺🇸 USD: %.4f\n"+
					"🇷🇺 RUB: %.4f\n"+
					"🇪🇺 EUR: %.4f\n"+
					"🇨🇳 CNY: %.4f",
				currentDate,
				ratesMap["USD"],
				ratesMap["RUB"],
				ratesMap["EUR"],
				ratesMap["CNY"],
			)
		}

	default:
		// Пытаемся понять, это ли расход: сумма категория
		parts := strings.Fields(text)
		if len(parts) == 2 {
			amount, err := strconv.ParseFloat(parts[0], 64)
			if err == nil {
				category := parts[1]
				err := AddExpense(userID, amount, category)
				if err != nil {
					msg.Text = "Ошибка при добавлении расхода"
				} else {
					msg.Text = fmt.Sprintf("Добавлено %.2f в категорию '%s'", amount, category)
				}
			} else {
				msg.Text = "Неверный формат. Пример: 200 еда"
			}
		} else {
			msg.Text = "Выберите действие из меню 👇"
		}
	}

	bot.Send(msg)
}
