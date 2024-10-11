package main

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"
)

// Функция для сбора наград
func collectBounties(carpet *Carpet, bounties []Bounty, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, bounty := range bounties {
		if math.Abs(carpet.X-bounty.X) < 1 && math.Abs(carpet.Y-bounty.Y) < 1 {
			// Логика сбора награды
			fmt.Printf("Carpet %s collected bounty at (%f, %f)\n", carpet.ID, bounty.X, bounty.Y)
			// Увеличиваем количество золота или здоровья
		} else {
			moveTowards(carpet, bounty.X, bounty.Y)
		}
	}
}

// Функция для уклонения от аномалий
func avoidAnomalies(carpet *Carpet, anomalies []Anomaly, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, anomaly := range anomalies {
		distance := math.Sqrt(math.Pow(carpet.X-anomaly.X, 2) + math.Pow(carpet.Y-anomaly.Y, 2))
		if distance < anomaly.Radius {
			// Уклоняемся от аномалии
			fmt.Printf("Carpet %s avoiding anomaly at (%f, %f)\n", carpet.ID, anomaly.X, anomaly.Y)
			moveAwayFrom(carpet, anomaly.X, anomaly.Y)
		}
	}
}

// Функция для атаки на соперников
func attackEnemies(carpet *Carpet, enemies []Carpet, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, enemy := range enemies {
		if math.Abs(carpet.X-enemy.X) < 5 && math.Abs(carpet.Y-enemy.Y) < 5 { // В пределах радиуса атаки
			// Логика атаки
			fmt.Printf("Carpet %s attacking enemy %s\n", carpet.ID, enemy.ID)
			// Например, отправляем API-запрос для атаки
		}
	}
}

// Функция для защиты
func defend(carpet *Carpet, healthThreshold int, wg *sync.WaitGroup) {
	defer wg.Done()
	if carpet.Health < healthThreshold {
		// Активируем щит
		carpet.Shield = true
		fmt.Printf("Carpet %s activated shield\n", carpet.ID)
	}
}

// Функция для управления ковром
func manageCarpet(carpet *Carpet, bounties []Bounty, anomalies []Anomaly, enemies []Carpet, wg *sync.WaitGroup) {
	// Запускаем горутины для каждой задачи
	wg.Add(4)
	go collectBounties(carpet, bounties, wg)
	go avoidAnomalies(carpet, anomalies, wg)
	go attackEnemies(carpet, enemies, wg)
	go defend(carpet, 30, wg) // Порог здоровья 30
}

// Функция для перемещения к цели
func moveTowards(carpet *Carpet, targetX, targetY float64) {
	dx := targetX - carpet.X
	dy := targetY - carpet.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 0 {
		dx /= distance
		dy /= distance
		carpet.X += dx * 0.1 // Умножаем на скорость
		carpet.Y += dy * 0.1 // Умножаем на скорость
	}
}

// Функция для уклонения от цели
func moveAwayFrom(carpet *Carpet, targetX, targetY float64) {
	dx := carpet.X - targetX
	dy := carpet.Y - targetY
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 0 {
		dx /= distance
		dy /= distance
		carpet.X += dx * 0.1 // Умножаем на скорость
		carpet.Y += dy * 0.1 // Умножаем на скорость
	}
}

// Функция для запроса состояния игры
func getGameState() (*GameState, error) {
	resp, err := apiRequest("GET", "/play/magcarp/player/move", nil) // Реализуйте apiRequest
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var state GameState
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return nil, err
	}

	return &state, nil
}

// Функция для обработки состояния игры
func processGameState(state *GameState, wg *sync.WaitGroup) {
	defer wg.Done()
	var wgCarpet sync.WaitGroup

	for _, carpet := range state.Carpets {
		wgCarpet.Add(4)
		go collectBounties(&carpet, state.Bounties, &wgCarpet)
		go avoidAnomalies(&carpet, state.Anomalies, &wgCarpet)
		go attackEnemies(&carpet, state.Carpets, &wgCarpet)
		go defend(&carpet, 30, &wgCarpet) // Порог здоровья 30
	}

	wgCarpet.Wait() // Ждем завершения всех горутин для ковров
}

// Главная функция бота
func runBot() {
	ticker := time.NewTicker(333 * time.Millisecond) // Устанавливаем интервал на 3 запроса в секунду
	defer ticker.Stop()

	var wg sync.WaitGroup

	for {
		select {
		case <-ticker.C:
			wg.Add(1)
			state, err := getGameState() // Получаем состояние игры
			if err != nil {
				fmt.Println("Error getting game state:", err)
				continue // Пропускаем итерацию при ошибке
			}
			go processGameState(state, &wg) // Обрабатываем состояние игры в горутине
		}
	}

	wg.Wait() // Ждем завершения всех обработок перед выходом
}
