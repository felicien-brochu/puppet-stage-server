package model

import "github.com/google/uuid"

// Board describes a board of multiple Devices (servos etc.)
type Board struct {
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	Servos  map[string]Servo `json:"servos"`
	BusType BusType          `json:"busType"`
}

// BusType type of bus
type BusType string

const (
	// BusTypeSerial serial bus type
	BusTypeSerial = BusType("serial")

	// BusTypeTCP TCP bus type
	BusTypeTCP = BusType("tcp")
)

// DefaultBoard returns a default board with a new ID
func DefaultBoard() Board {
	return Board{
		uuid.New().String(),
		"",
		make(map[string]Servo),
		BusTypeSerial,
	}
}
