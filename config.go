package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	TelegramBotToken   string
	PostgresConnString string
)

func init() {
	// Загружаем .env только в локальной среде
	env := os.Getenv("ENV")
	if env != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, using system environment variables")
		}
	}

	TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	if TelegramBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	PostgresConnString = os.Getenv("DATABASE_URL")
	if PostgresConnString == "" {
		log.Fatal("DATABASE_URL is not set")
	}
}
