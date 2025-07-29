package main

import (
    "os"
)

const (
	TelegramBotToken = "8028673613:AAGXhTldnMKzEOaCaTm_eCqCuDWyygxzMPk"
	//WeatherAPIKey    = "ВАШ_OPENWEATHERMAP_API_KEY"
	//PostgresConnString = "host=localhost port=5432 user=app password=pass dbname=expensesbot sslmode=disable"
)

//package main

var PostgresConnString string

func init() {
    PostgresConnString = os.Getenv("DATABASE_URL")
    if PostgresConnString == "" {
        // fallback на локальный конфиг для разработки
        //PostgresConnString = "host=localhost port=5432 user=adm password=adm dbname=bot sslmode=disable"
		PostgresConnString = "host=localhost port=5432 user=app password=pass dbname=expensesbot sslmode=disable"
    }
}

