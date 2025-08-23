package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"
	"net/http"
	"strings"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type ValCurs struct {
	XMLName xml.Name  `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID       string `xml:"Id,attr"`
	CharCode string  `xml:"CharCode"`
	Nominal  int     `xml:"Nominal"`
	Name  	 string  `xml:"Name"`
	Value    string  `xml:"Value"` // Значение как строка (с запятой)
}

func GetExchangeRatesFromXML() (string, error) {
	currentDate := time.Now().Format("2006-01-02")
	url := fmt.Sprintf("https://nbt.tj/tj/kurs/export_xml.php?date=%s&export=xmlout", currentDate)
	
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	// Декодируем windows-1251 → utf-8
	decoded := transform.NewReader(resp.Body, charmap.Windows1251.NewDecoder())

	// Читаем всё содержимое
	xmlBytes, err := io.ReadAll(decoded)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения XML: %w", err)
	}

	// Заменяем windows-1251 на utf-8 в заголовке
	xmlStr := strings.Replace(string(xmlBytes), `encoding="windows-1251"`, `encoding="utf-8"`, 1)

	// Парсим
	var valCurs ValCurs
	err = xml.Unmarshal([]byte(xmlStr), &valCurs)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга XML: %w", err)
	}
	
	// Ищем курсы
	var usd, rub, eur, cny string
	for _, record := range valCurs.Valutes {
		code := strings.TrimSpace(record.CharCode)
		switch code {
		case "USD":
			usd = record.Value
		case "RUB":
			rub = record.Value
		case "EUR":
			eur = record.Value
		case "CNY":
			cny = record.Value
		}
	}
	
	if usd == "" || rub == "" || eur == "" || cny == "" {
		return "", fmt.Errorf("не найдены обе валюты")
	}

	return fmt.Sprintf("USD: %s, RUB: %s, EUR: %s, CNY: %s", usd, rub, eur, cny), nil
}

