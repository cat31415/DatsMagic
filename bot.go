package main

import (
	"fmt"
	"math"
	"sync"
)

// Функция для обновления состояния всех ковров
func updateCarpets(carpets []Carpet, accelerations []Vector, attacks []Vector) (*MoveResponse, error) {
	req := PlayerCommand{Transports: []TransportCommand{}}

	for i, carpet := range carpets {
		// Добавляем команду для ковра с учетом ускорения и атаки
		req.Transports = append(req.Transports, TransportCommand{
			ID:           carpet.ID,
			Acceleration: accelerations[i],
		})
		if attacks[i].X != 0 || attacks[i].Y != 0 {
			req.Transports[i].Attack = attacks[i]
		}
	}

	// Отправляем команду на сервер
	return sendPlayerCommand(req)
}

// Функция для атаки на соперников
func attackEnemies(carpet *Carpet, enemies []Carpet, attackRadius float64, wg *sync.WaitGroup, attackResults chan Vector) {
	defer wg.Done()
	attackVector := Vector{X: 0, Y: 0} // По умолчанию атака не происходит

	var closestEnemy *Carpet
	var minDistance float64 = attackRadius + 1 // Минимальная дистанция для выбора цели

	// Поиск ближайшего врага в пределах радиуса атаки
	for _, enemy := range enemies {
		// Рассчитываем расстояние до врага
		dx := carpet.X - enemy.X
		dy := carpet.Y - enemy.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		// Проверяем, находится ли враг в радиусе атаки
		if distance != 0 && distance <= attackRadius && distance < minDistance {
			minDistance = distance
			closestEnemy = &enemy
		}
	}

	// Если есть ближайший враг, атакуем
	if closestEnemy != nil {
		fmt.Printf("Carpet %s attacking enemy %s at distance %.2f\n", carpet.ID, closestEnemy.ID, minDistance)

		// Рассчитываем вектор атаки в направлении к врагу
		dx := closestEnemy.X - carpet.X
		dy := closestEnemy.Y - carpet.Y
		attackVector = Vector{X: dx, Y: dy}

		// Нормализуем вектор атаки (чтобы не зависеть от расстояния)
		distance := math.Sqrt(dx*dx + dy*dy)
		if distance > 0 {
			attackVector.X /= distance
			attackVector.Y /= distance
		}
	}

	// Отправляем вектор атаки в канал
	attackResults <- attackVector
}

// Функция для расчета ускорения на основе целевой точки
func calculateAcceleration(carpet *Carpet, targetX, targetY float64) Vector {
	dx := targetX - carpet.X
	dy := targetY - carpet.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 0 {
		// Нормализуем вектор и применяем максимальное ускорение
		accel := carpet.MaxAccel
		dx /= distance
		dy /= distance
		return Vector{X: dx * accel, Y: dy * accel}
	}
	return Vector{X: 0, Y: 0} // Если находимся на месте, возвращаем нулевое ускорение
}

// Функция для уклонения от аномалий
func avoidAnomalies(carpet *Carpet, anomalies []Anomaly, wg *sync.WaitGroup, results chan Vector) {
	defer wg.Done()
	var totalAnomalyAcceleration Vector

	for _, anomaly := range anomalies {
		distance := math.Sqrt(math.Pow(carpet.X-anomaly.X, 2) + math.Pow(carpet.Y-anomaly.Y, 2))
		if distance < anomaly.Radius {
			fmt.Printf("Carpet %s avoiding anomaly at (%f, %f)\n", carpet.ID, anomaly.X, anomaly.Y)
			// Рассчитываем ускорение для уклонения
			acceleration := moveAwayFrom(carpet, anomaly.X, anomaly.Y)
			totalAnomalyAcceleration.X += acceleration.X
			totalAnomalyAcceleration.Y += acceleration.Y
		}
	}

	results <- totalAnomalyAcceleration // Отправляем результирующее ускорение в канал
}

// Функция для сбора наград
func collectBounties(carpet *Carpet, bounties []Bounty, wg *sync.WaitGroup, results chan Vector) {
	defer wg.Done()
	var totalSelfAcceleration Vector

	for _, bounty := range bounties {
		if math.Abs(carpet.X-bounty.X) < 1 && math.Abs(carpet.Y-bounty.Y) < 1 {
			fmt.Printf("Carpet %s collected bounty at (%f, %f)\n", carpet.ID, bounty.X, bounty.Y)
		} else {
			acceleration := calculateAcceleration(carpet, bounty.X, bounty.Y)
			totalSelfAcceleration.X += acceleration.X
			totalSelfAcceleration.Y += acceleration.Y
		}
	}
	fmt.Println(totalSelfAcceleration)
	if carpet.SelfAcceleration != totalSelfAcceleration {
		fmt.Print("Changed selfAcceleration")
	}
	results <- totalSelfAcceleration // Отправляем результирующее ускорение в канал
}

// Функция для уклонения от цели
func moveAwayFrom(carpet *Carpet, targetX, targetY float64) Vector {
	dx := carpet.X - targetX
	dy := carpet.Y - targetY
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 0 {
		// Нормализуем вектор направления
		dx /= distance
		dy /= distance

		// Рассчитываем новое ускорение для уклонения
		acceleration := Vector{
			X: dx * carpet.MaxAccel, // Умножаем на максимальное ускорение
			Y: dy * carpet.MaxAccel, // Умножаем на максимальное ускорение
		}
		return acceleration
	}
	return Vector{X: 0, Y: 0} // Если ковер на месте, возвращаем нулевое ускорение
}

// Функция для управления коврами
func manageCarpet(carpet *Carpet, bounties []Bounty, anomalies []Anomaly, enemies []Carpet, wg *sync.WaitGroup, results chan Vector, attackResults chan Vector) {
	wg.Add(3) // Увеличиваем счетчик горутин для атаки

	// Горутин для сбора наград
	go collectBounties(carpet, bounties, wg, results)

	// Горутин для уклонения от аномалий
	go avoidAnomalies(carpet, anomalies, wg, results)

	// Горутин для атаки на соперников
	go attackEnemies(carpet, enemies, 30.0, wg, attackResults) // Радиус атаки, например, 5 метров
}

// Основная функция бота
func runBot(state *MoveResponse) *MoveResponse {

	var wg sync.WaitGroup
	results := make(chan Vector, len(state.Carpets))       // Канал для ускорений
	attackResults := make(chan Vector, len(state.Carpets)) // Канал для векторов атак

	for _, carpet := range state.Carpets {
		carpet.MaxAccel = state.MaxAccel
		manageCarpet(&carpet, state.Bounties, state.Anomalies, state.Carpets, &wg, results, attackResults)
	}

	wg.Wait()            // Ждем завершения всех горутин
	close(results)       // Закрываем канал для ускорений
	close(attackResults) // Закрываем канал для атак

	// Сбор итоговых результатов
	accelerations := make([]Vector, len(state.Carpets))
	attacks := make([]Vector, len(state.Carpets)) // Список для хранения векторов атак

	for i := range state.Carpets {
		accelerations[i] = <-results
		attacks[i] = <-attackResults
	}
	// Передаем вектора атак и ускорений в функцию sendPlayerCommand
	res, err := updateCarpets(state.Carpets, accelerations, attacks)
	if err != nil {
		fmt.Printf("Error updating carpets: %s\n", err)
	}
	return res
}
