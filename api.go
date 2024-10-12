package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Функция для отправки команды на действие
func sendPlayerCommand(command PlayerCommand) (*MoveResponse, error) {
	url := "https://games-test.datsteam.dev/play/magcarp/player/move"

	jsonData, err := json.Marshal(command)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при кодировании JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("Ошибка при создании запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", "66fbdaf5594c466fbdaf5594c8") // замените на ваш актуальный токен

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ошибка: статус код %d", resp.StatusCode)
	}

	var moveResponse MoveResponse
	if err := json.NewDecoder(resp.Body).Decode(&moveResponse); err != nil {
		return nil, fmt.Errorf("Ошибка при декодировании ответа: %v", err)
	}

	return &moveResponse, nil
}

// Функция для получения информации по карте и коврам
func getInitialState() (*MoveResponse, error) {
	url := "https://games-test.datsteam.dev/play/magcarp/player/move" // заменить на актуальный URL

	// Создаем тело запроса с пустым массивом transports
	playerCommand := PlayerCommand{
		Transports: []TransportCommand{},
	}

	// Кодируем тело запроса в JSON
	jsonData, err := json.Marshal(playerCommand)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при кодировании JSON: %v", err)
	}

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("Ошибка при создании запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", "66fbdaf5594c466fbdaf5594c8") // замените на актуальный токен

	client := &http.Client{Timeout: time.Second / 3}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ошибка: статус код %d", resp.StatusCode)
	}

	// Декодируем JSON-ответ в структуру MoveResponse
	var moveResponse MoveResponse
	if err := json.NewDecoder(resp.Body).Decode(&moveResponse); err != nil {
		return nil, fmt.Errorf("Ошибка при декодировании ответа: %v", err)
	}

	return &moveResponse, nil
}
