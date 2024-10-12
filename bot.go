package main

import (
	"fmt"
	"math"
	"time"
)

// Функция для обновления состояния всех ковров
func updateCarpets(carpets []Carpet, accelerations []Vector, attacks []Vector, shields []bool) (*MoveResponse, error) {
	req := PlayerCommand{Transports: []TransportCommand{}}

	for i, carpet := range carpets {
		// Добавляем команду для ковра с учетом ускорения и атаки
		req.Transports = append(req.Transports, TransportCommand{
			ID:             carpet.ID,
			Acceleration:   accelerations[i],
			ActivateShield: shields[i], // Указываем, нужно ли активировать щит
		})

		// Если есть вектор атаки, добавляем его
		if attacks[i].X != 0 || attacks[i].Y != 0 {

			req.Transports[i].Attack = attacks[i]
		}
	}

	// Отправляем команду на сервер
	return sendPlayerCommand(req)
}

// Функция для атаки на соперников с предсказанием и активацией щита
func attackEnemies(carpet *Carpet, enemies []Enemy, attackRadius float64) (Vector, bool) {
	attackVector := Vector{X: 0, Y: 0} // По умолчанию атака не происходит
	activateShield := false            // По умолчанию щит не активируется

	var weakestEnemy *Enemy
	minHealth := 200

	// Поиск врага с самым низким здоровьем в пределах радиуса атаки
	for _, enemy := range enemies {
		if enemy.Status != "alive" {
			continue // Игнорируем мертвых врагов
		}

		// Рассчитываем расстояние до врага
		dx := carpet.X - enemy.X
		dy := carpet.Y - enemy.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		// Проверяем, находится ли враг в радиусе атаки
		if distance <= attackRadius {
			// Если здоровье врага меньше текущего минимального
			if enemy.Health < minHealth {
				minHealth = enemy.Health
				weakestEnemy = &enemy // Запоминаем врага с минимальным здоровьем
			}
		} else {
			// Предсказание будущей позиции врага через 1 секунду
			predictedEnemyX := enemy.X + enemy.Velocity.X*1 // Позиция через 1 секунду
			predictedEnemyY := enemy.Y + enemy.Velocity.Y*1

			// Предсказание вашей позиции через 1 секунду
			predictedCarpetX := carpet.X + carpet.Velocity.X*1 + carpet.MaxAccel*1 // Позиция через 1 секунду
			predictedCarpetY := carpet.Y + carpet.Velocity.Y*1 + carpet.MaxAccel*1

			// Рассчитываем расстояние до предсказанной позиции врага
			predictedDx := predictedCarpetX - predictedEnemyX
			predictedDy := predictedCarpetY - predictedEnemyY
			predictedDistance := math.Sqrt(predictedDx*predictedDx + predictedDy*predictedDy)

			// Если предсказанная дистанция меньше радиуса атаки, активируем щит
			if predictedDistance <= attackRadius {
				fmt.Println("Shield active")
				activateShield = true
			}
		}
	}

	// Если есть враг с минимальным здоровьем, атакуем
	if weakestEnemy != nil {
		//fmt.Printf("Carpet %s attacking enemy with health %d at distance %.2f\n", carpet.ID, weakestEnemy.Health, minHealth)

		attackVector = Vector{X: weakestEnemy.X, Y: weakestEnemy.Y}

	}

	return attackVector, activateShield
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

// Функция для уклонения от аномалий с учетом их скорости и минимального изменения маршрута
func avoidAnomalies(carpet *Carpet, anomalies []Anomaly) Vector {
	totalAvoidance := Vector{X: 0, Y: 0} // Вектор уклонения

	for _, anomaly := range anomalies {
		// Рассчитываем текущее расстояние до аномалии
		dx := carpet.X - anomaly.X
		dy := carpet.Y - anomaly.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		// Если аномалия слишком близка, начнем избегать её
		if distance <= anomaly.Radius+2 { // Например, начинаем избегать за 5 метров до радиуса
			// Предсказание будущей позиции аномалии через 1 секунду
			predictedAnomalyX := anomaly.X + anomaly.Velocity.X*1
			predictedAnomalyY := anomaly.Y + anomaly.Velocity.Y*1

			// Рассчитываем направление к аномалии
			predictedDx := carpet.X - predictedAnomalyX
			predictedDy := carpet.Y - predictedAnomalyY
			predictedDistance := math.Sqrt(predictedDx*predictedDx + predictedDy*predictedDy)

			// Если предсказанное расстояние слишком маленькое, нужно уклониться
			if predictedDistance < anomaly.Radius {
				// Рассчитываем вектор уклонения "вбок"
				sideStepX := -predictedDy // Сдвиг вбок относительно направления аномалии
				sideStepY := predictedDx  // Это создает перпендикулярный вектор

				// Нормализуем вектор уклонения
				sideStepDistance := math.Sqrt(sideStepX*sideStepX + sideStepY*sideStepY)
				if sideStepDistance > 0 {
					sideStepX /= sideStepDistance
					sideStepY /= sideStepDistance
				}

				// Добавляем вектор уклонения с минимальным изменением маршрута
				avoidanceStrength := 0.5 // Степень отклонения от маршрута
				totalAvoidance.X += sideStepX * avoidanceStrength
				totalAvoidance.Y += sideStepY * avoidanceStrength
			}
		}
	}

	return totalAvoidance
}

// Функция для сбора наград
func collectBounties(carpet *Carpet, bounties []Bounty) Vector {
	var totalSelfAcceleration Vector
	var closestBounty *Bounty
	minDistance := math.MaxFloat64 // Начальное значение для минимального расстояния

	// Поиск ближайшей монеты
	for _, bounty := range bounties {
		dx := carpet.X - bounty.X
		dy := carpet.Y - bounty.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		// Проверяем, есть ли монета ближе, чем предыдущая
		if distance < minDistance {
			minDistance = distance
			closestBounty = &bounty // Сохраняем ближайшую монету
		}
	}

	// Если нашли ближайшую монету, рассчитываем ускорение
	if closestBounty != nil {
		// Проверяем, достигли ли мы монеты
		if minDistance < 1 { // Порог для "собрания" монеты
			fmt.Printf("Carpet %s collected bounty at (%f, %f)\n", carpet.ID, closestBounty.X, closestBounty.Y)
		} else {
			// Рассчитываем желаемое ускорение к ближайшей монете
			desiredAcceleration := calculateAcceleration(carpet, closestBounty.X, closestBounty.Y)

			// Рассчитываем текущее направление
			currentVelocity := Vector{X: carpet.Velocity.X, Y: carpet.Velocity.Y}
			currentSpeed := math.Sqrt(currentVelocity.X*currentVelocity.X + currentVelocity.Y*currentVelocity.Y)

			// Нормализуем текущее направление
			if currentSpeed > 0 {
				direction := Vector{X: currentVelocity.X / currentSpeed, Y: currentVelocity.Y / currentSpeed}

				// Отклоняем желаемое ускорение в сторону текущего направления
				adjustedAcceleration := Vector{
					X: desiredAcceleration.X - direction.X*(desiredAcceleration.X*direction.X+desiredAcceleration.Y*direction.Y),
					Y: desiredAcceleration.Y - direction.Y*(desiredAcceleration.X*direction.X+desiredAcceleration.Y*direction.Y),
				}

				// Добавляем скорректированное ускорение к общему ускорению
				totalSelfAcceleration.X += adjustedAcceleration.X
				totalSelfAcceleration.Y += adjustedAcceleration.Y
			} else {
				// Если ковер стоит на месте, просто добавляем желаемое ускорение
				totalSelfAcceleration.X += desiredAcceleration.X
				totalSelfAcceleration.Y += desiredAcceleration.Y
			}
		}
	}

	return totalSelfAcceleration
}

// Функция для управления коврами
func manageCarpet(carpet *Carpet, bounties []Bounty, anomalies []Anomaly, enemies []Enemy, attackRange int) (Vector, Vector, bool) {

	// Сбор наград
	bountyVector := collectBounties(carpet, bounties)

	// Уклонение от аномалий
	//avoidanceVector := avoidAnomalies(carpet, anomalies)

	// Атака на врагов и активация щита
	attackVector, activateShield := attackEnemies(carpet, enemies, float64(attackRange))
	var finalAcceleration Vector
	// if carpet.AnomalyAcceleration.X != 0 || carpet.AnomalyAcceleration.Y != 0 {
	// 	fmt.Println("Carpet", carpet.ID, "is using anomaly acceleration")
	// 	finalAcceleration = Vector{
	// 		X: avoidanceVector.X,
	// 		Y: avoidanceVector.Y,
	// 	}
	// } else {

	// }
	finalAcceleration = Vector{
		X: bountyVector.X,
		Y: bountyVector.Y,
	}
	// Рассчитываем длину итогового вектора
	length := math.Sqrt(finalAcceleration.X*finalAcceleration.X + finalAcceleration.Y*finalAcceleration.Y)

	// Если длина вектора больше максимального ускорения, нормализуем вектор
	if length > carpet.MaxAccel && length > 0 {
		finalAcceleration.X = (finalAcceleration.X / length) * carpet.MaxAccel
		finalAcceleration.Y = (finalAcceleration.Y / length) * carpet.MaxAccel
	}

	return finalAcceleration, attackVector, activateShield
}

// Основная функция бота
func runBot(state *MoveResponse) *MoveResponse {
	if state == nil {
		fmt.Println("state is nil")
		var err error
		time.Sleep(time.Second / 3)
		state, err = getInitialState()
		if err != nil {
			fmt.Printf("Error getting initial state: %s\n", err)
		}
	}
	accelerations := make([]Vector, len(state.Carpets))
	attacks := make([]Vector, len(state.Carpets))
	shields := make([]bool, len(state.Carpets))

	for i, carpet := range state.Carpets {
		carpet.MaxAccel = state.MaxAccel

		// Управление ковром
		acceleration, attackVector, activateShield := manageCarpet(&carpet, state.Bounties, state.Anomalies, state.Enemies, 200)
		accelerations[i] = acceleration
		attacks[i] = attackVector
		shields[i] = activateShield
	}

	// Передаем вектора атак и ускорений в функцию sendPlayerCommand
	req, err := updateCarpets(state.Carpets, accelerations, attacks, shields)
	if err != nil {
		fmt.Printf("Error updating carpets: %s\n", err)
	}
	return req
}
