package main

import (
	"log"
	"os"
)

var (
	TelegramBotToken string
	PostgresConnString string
)

func init() {
	TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	log.Println("TELEGRAM_BOT_TOKEN = ", TelegramBotToken)
    if TelegramBotToken == "" {
        log.Fatal("TELEGRAM_BOT_TOKEN is not set")
    }

    PostgresConnString = os.Getenv("DATABASE_URL")
    if PostgresConnString == "" {
		log.Fatal("DATABASE_URL is not set")
    }
	log.Printf("Config loaded: TelegramBotToken=*****, DB=****")
}

