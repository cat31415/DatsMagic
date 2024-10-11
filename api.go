package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Структуры для отправки запроса
type PlayerCommand struct {
	Transports []TransportCommand `json:"transports"`
}

type TransportCommand struct {
	ID             string `json:"id,omitempty"`
	Acceleration   Vector `json:"acceleration,omitempty"`
	ActivateShield bool   `json:"activateShield,omitempty"`
	Attack         Vector `json:"attack,omitempty"`
}

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Структура для получения ответа (здесь добавь поля, которые ожидаешь получить)
type MoveResponse struct {
	Name       string      `json:"name"`
	X          int         `json:"x"`
	Y          int         `json:"y"`
	Health     int         `json:"health"`
	Transports []Transport `json:"transports"`
	Anomalies  []Anomaly   `json:"anomalies"`
	Bounties   []Bounty    `json:"bounties"`
	Enemies    []Enemy     `json:"enemies"`
	MapSize    Coordinate  `json:"mapSize"`
}

type Transport struct {
	ID     string `json:"id"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Health int    `json:"health"`
}

type Anomaly struct {
	ID     string  `json:"id"`
	X      int     `json:"x"`
	Y      int     `json:"y"`
	Radius float64 `json:"radius"`
}

type Bounty struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Points int `json:"points"`
}

type Enemy struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Health int    `json:"health"`
	Status string `json:"status"`
}

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
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

	client := &http.Client{Timeout: 10 * time.Second}
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

func main() {
	// Получаем начальное состояние карты и ковров
	state, err := getInitialState()
	if err != nil {
		log.Fatalf("Ошибка при получении начального состояния: %v", err)
	}

	fmt.Printf("Игрок %s находится на координатах X=%d, Y=%d\n", state.Name, state.X, state.Y)
	fmt.Printf("Здоровье игрока: %d\n", state.Health)
	fmt.Println("Доступные транспорты:")
	for _, transport := range state.Transports {
		fmt.Printf("Транспорт ID=%s на координатах X=%d, Y=%d, здоровье: %d\n", transport.ID, transport.X, transport.Y, transport.Health)
	}
	fmt.Println("Аномалии на карте:")
	for _, anomaly := range state.Anomalies {
		fmt.Printf("Аномалия ID=%s на координатах X=%d, Y=%d, радиус: %.2f\n", anomaly.ID, anomaly.X, anomaly.Y, anomaly.Radius)
	}
	fmt.Println("Награды на карте:")
	for _, bounty := range state.Bounties {
		fmt.Printf("Награда на координатах X=%d, Y=%d, очки: %d\n", bounty.X, bounty.Y, bounty.Points)
	}
	fmt.Println("Информация о врагах:")
	for _, enemy := range state.Enemies {
		fmt.Printf("Враг на координатах X=%d, Y=%d, здоровье: %d, статус: %s\n", enemy.X, enemy.Y, enemy.Health, enemy.Status)
	}
}
