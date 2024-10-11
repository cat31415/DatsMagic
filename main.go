package main

import (
	"fmt"
	"log"
)

func main() {
	// Получаем начальное состояние карты и ковров
	state, err := getInitialState()
	if err != nil {
		log.Fatalf("Ошибка при получении начального состояния: %v", err)
	}

	// Сохраняем состояние в файл "game_state.json"
	if err := saveStateToJSON(state, "game_state.json"); err != nil {
		log.Fatalf("Ошибка при сохранении состояния: %v", err)
	}

	fmt.Printf("Игрок %s находится на координатах X=%d, Y=%d\n", state.Name, state.X, state.Y)
	fmt.Printf("Здоровье игрока: %d\n", state.Health)
	fmt.Println("Доступные транспорты:")
	for _, transport := range state.Carpets {
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
