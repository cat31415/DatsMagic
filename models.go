package main

// Structure to represent the bot's carpet
type Carpet struct {
	ID                  string  `json:"id"`
	X                   float64 `json:"x"`
	Y                   float64 `json:"y"`
	Health              int     `json:"health"`
	MaxAccel            float64
	Velocity            Vector `json:"velocity"`
	AnomalyAcceleration Vector `json:"anomalyAcceleration"`
	SelfAcceleration    Vector `json:"selfAcceleration"`
	Shield              bool   `json:"activateShield"`
}

// Structure for vector representation
type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Structure for anomalies
type Anomaly struct {
	ID       string  `json:"id"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Radius   float64 `json:"effectiveRadius"`
	Strength float64 `json:"strength"`
	Velocity Vector  `json:"velocity"`
}

// Structure for bounties
type Bounty struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Points int     `json:"points"`
}

// Structure for enemies
type Enemy struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Health   int     `json:"health"`
	Status   string  `json:"status"`
	Velocity Vector  `json:"velocity"`
}

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

// Структура для получения ответа (здесь добавь поля, которые ожидаешь получить)
type MoveResponse struct {
	Name                  string     `json:"name"`
	X                     int        `json:"x"`
	Y                     int        `json:"y"`
	Health                int        `json:"health"`
	Carpets               []Carpet   `json:"transports"`
	Anomalies             []Anomaly  `json:"anomalies"`
	Bounties              []Bounty   `json:"bounties"`
	Enemies               []Enemy    `json:"enemies"`
	MapSize               Coordinate `json:"mapSize"`
	MaxAccel              float64    `json:"maxAccel"`
	AttackExplosionRadius float64    `json:"attackExplosionRadius"`
	AttackRange           float64    `json:"attackRange"`
	ReviveTimeoutSec      int        `json:"reviveTimeoutSec"`
}

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}
