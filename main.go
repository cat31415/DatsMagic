package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Second / 3) // создаем тикер для выполнения запроса 3 раза в секунду
	defer ticker.Stop()

	counter := 1
	state, err := sendPlayerCommand(PlayerCommand{[]TransportCommand{}})
	if err != nil {
		fmt.Println("err 429")
	}

	interval := time.Second * 3
	// Таймер для отсчета времени
	maxRequests := 8
	requestTicker := time.NewTicker(375 * time.Millisecond) // Выполняем запросы каждые 0.375 секунды (8 за 3 секунды)
	intervalTicker := time.NewTicker(interval)
	defer requestTicker.Stop()
	defer intervalTicker.Stop()
	requests := 1

	for {
		select {
		case <-requestTicker.C:
			if requests < maxRequests {
				// Пример команды для одного из ковров
				state = runBot(state)

				// Сохраняем команду и ответ в JSON
				commandFilename := fmt.Sprintf("/comands/command_%d.json", counter)
				if err := saveToJSON(commandFilename, state); err != nil {
					log.Printf("Ошибка при сохранении команды в JSON: %v", err)
				}
				requests++
				counter++
			}
		case <-intervalTicker.C:
			// Сброс количества запросов через каждые 3 секунды
			requests = 0
		}
	}
}
