package model

// Board describes a board of multiple Devices (servos etc.)
type Board struct {
	Name   string
	Servos []*Servo
	Bus    Bus
}
