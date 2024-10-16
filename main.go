package main

import (
	"fmt"
	"time"
)

func main() {

	counter := 1
	state, err := sendPlayerCommand(PlayerCommand{[]TransportCommand{}})
	if err != nil {
		fmt.Println("err 429")
	}
	time.Sleep(time.Second / 2)
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
