package model

// Board describes a board of multiple Devices (servos etc.)
type Board struct {
	Name    string   `json:"name"`
	Servos  []*Servo `json:"servos"`
	Bus     Bus      `json:"-"`
	BusType BusType  `json:"busType"`
}
