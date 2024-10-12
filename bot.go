package main

import (
	"fmt"
	"math"
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
func attackEnemies(carpet *Carpet, enemies []Carpet, attackRadius float64) (Vector, bool) {
	attackVector := Vector{X: 0, Y: 0} // По умолчанию атака не происходит
	activateShield := false            // По умолчанию щит не активируется

	var closestEnemy *Carpet
	var minDistance float64 = attackRadius + 1 // Минимальная дистанция для выбора цели

	// Поиск ближайшего врага в пределах радиуса атаки
	for _, enemy := range enemies {
		// Рассчитываем расстояние до врага
		dx := carpet.X - enemy.X
		dy := carpet.Y - enemy.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		// Проверяем, находится ли враг в радиусе атаки
		if distance != 0 {
			if distance <= attackRadius && distance < minDistance {
				minDistance = distance
				closestEnemy = &enemy
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
					fmt.Println("Sheilde active")
					activateShield = true
				}
			}
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

	for _, bounty := range bounties {
		if math.Abs(carpet.X-bounty.X) < 1 && math.Abs(carpet.Y-bounty.Y) < 1 {
			fmt.Printf("Carpet %s collected bounty at (%f, %f)\n", carpet.ID, bounty.X, bounty.Y)
		} else {
			acceleration := calculateAcceleration(carpet, bounty.X, bounty.Y)
			totalSelfAcceleration.X += acceleration.X
			totalSelfAcceleration.Y += acceleration.Y
		}
	}

	return totalSelfAcceleration
}

// Функция для управления коврами
func manageCarpet(carpet *Carpet, bounties []Bounty, anomalies []Anomaly, enemies []Carpet) (Vector, Vector, bool) {
	// Сбор наград
	bountyVector := collectBounties(carpet, bounties)

	// Уклонение от аномалий
	avoidanceVector := avoidAnomalies(carpet, anomalies)

	// Атака на врагов и активация щита
	attackVector, activateShield := attackEnemies(carpet, enemies, 30.0)

	// Суммируем оба вектора
	finalAcceleration := Vector{
		X: bountyVector.X + avoidanceVector.X,
		Y: bountyVector.Y + avoidanceVector.Y,
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
	accelerations := make([]Vector, len(state.Carpets))
	attacks := make([]Vector, len(state.Carpets))
	shields := make([]bool, len(state.Carpets))

	for i, carpet := range state.Carpets {
		carpet.MaxAccel = state.MaxAccel

		// Управление ковром
		acceleration, attackVector, activateShield := manageCarpet(&carpet, state.Bounties, state.Anomalies, state.Carpets)
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
