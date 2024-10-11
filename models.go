package main

// Structure to represent the bot's carpet
type Carpet struct {
	ID                  string  `json:"id"`
	X                   float64 `json:"x"`
	Y                   float64 `json:"y"`
	Health              int     `json:"health"`
	Velocity            Vector  `json:"velocity"`
	AnomalyAcceleration Vector  `json:"anomalyAcceleration"`
	SelfAcceleration    Vector  `json:"selfAcceleration"`
	Shield              bool    `json:"shield"`
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
	Radius   float64 `json:"radius"`
	Strength float64 `json:"strength"`
}

// Structure for bounties
type Bounty struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Points int     `json:"points"`
}

// Structure for game state
type GameState struct {
	Carpets   []Carpet  `json:"carpets"`
	Anomalies []Anomaly `json:"anomalies"`
	Bounties  []Bounty  `json:"bounties"`
}
