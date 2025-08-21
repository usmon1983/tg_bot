package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shirou/gopsutil/v3/process"
)

func main() {
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

	// Запускаем мониторинг зависших процессов в отдельной горутине
	go monitorProcesses()

	//menu := tgbotapi.NewMessage(0, "Меню")
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("➕ Добавить расход"),
			tgbotapi.NewKeyboardButton("📊 Статистика"),
			tgbotapi.NewKeyboardButton("🌦️ Прогноз погоды"),
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

// мониторинг процессов
func monitorProcesses() {
	fmt.Println("Старт мониторинга процессов")

	selfPID := int32(os.Getpid())
	fmt.Printf("👤 Текущий PID: %d\n", selfPID)

	processes, err := process.Processes()
	if err != nil {
		fmt.Println("Ошибка получения процессов:", err)
		return
	}

	fmt.Printf("Найдено процессов: %d\n", len(processes))

	for _, p := range processes {

		ramSize, err := p.MemoryInfo()
		if err == nil && ramSize != nil {
			fmt.Printf("RSS: %d байт (%.2f МБ)\n", ramSize.RSS, float64(ramSize.RSS)/1024/1024)
		} else {
			if strings.Contains(err.Error(), "Access is denied") {
				fmt.Printf("Нет доступа к памяти процесса PID %d. Пропускаем.\n", p.Pid)
				return
			} else {
				fmt.Println("Ошибка MemoryInfo():", err)
			}
		}

		cpuPercent, err := p.CPUPercent()
		if err == nil {
			fmt.Printf("CPU: %.2f%%\n", cpuPercent)
		} else {
			if strings.Contains(err.Error(), "Access is denied") {
				fmt.Printf("❌ Нет доступа к CPU процессу PID %d. Пропускаем.\n", p.Pid)
				return
			} else {
				fmt.Println("❌ Ошибка CPUPercent():", err)
			}
		}

		isRunning, err := p.IsRunning()
		if err != nil {
			fmt.Println("Ошибка IsRunning():", err)
			if strings.Contains(err.Error(), "Access is denied") {
				fmt.Println("Нет доступа к информации о процессе. Возможно, это системный или защищённый процесс.")
			}
		} else {
			fmt.Println("Процесс активен:", isRunning)
		}

		name, err := p.Name()
		if err != nil {
			name = "неизвестно"
		}

		if p.Pid != selfPID {
			continue
		}

		const epsilon = 0.000001
		if cpuPercent < epsilon || !isRunning {
			if p != nil {
				p.Kill()
			}
		}

		fmt.Printf("Нашли БОТ-а PID: %d, Name: %s, RAM: %s, isRunning: %v\n", p.Pid, name, ramSize, isRunning)
		fmt.Printf("PID: %d, Name: %s\n", p.Pid, name)
	}
}
