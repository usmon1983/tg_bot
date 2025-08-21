package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func GetWheatherFromJSON(city string)  string{
	lang := "ru"
	fmt.Println(city)
	q := city
	
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Ошибка загрузки .env файла")
    }

	apiWheatherKey := os.Getenv("WHEATHER_API_KEY")
	urlApiWheather := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=%s", apiWheatherKey, q, lang)
    resp, err := http.Get(urlApiWheather)
	if err != nil {
        log.Fatal("Ошибка запроса:", err)
    }
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal("Ошибка чтения ответа:", err)
    }

    // Пример структуры для парсинга (можно уточнить)
    var data struct {
        Location struct {
            Name string `json:"name"`
        } `json:"location"`
        Current struct {
            TempC     float64 `json:"temp_c"`
            Condition struct {
                Text string `json:"text"`
            } `json:"condition"`
        } `json:"current"`
    }

    err = json.Unmarshal(body, &data)
    if err != nil {
        log.Fatal("Ошибка парсинга JSON:", err)
    }

    result := fmt.Sprintf("🌤️ Погода в %s: %s, %.1f°C",
        data.Location.Name, data.Current.Condition.Text, data.Current.TempC)

    return result
}