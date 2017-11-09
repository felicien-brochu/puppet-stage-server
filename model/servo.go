package model

import (
	"errors"

	"github.com/google/uuid"
)

// Servo describes a servomotor
type Servo struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Addr            int    `json:"addr"`
	DefaultPosition int    `json:"defaultPosition"`
	HardMin         int    `json:"hardMin"`
	HardMax         int    `json:"hardMax"`
	Min             int    `json:"min"`
	Max             int    `json:"max"`
	Inverted        bool   `json:"inverted"`
}

const (
	servoHardMin = 130
	servoHardMax = 600
	servoAvg     = (servoHardMin + servoHardMax) / 2
)

// IsOKValue checks if the given value is in bounds of the Device
func (servo *Servo) IsOKValue(value int) bool {
	return value >= servo.Min && value <= servo.Max
}

// DefaultServo constructs a new Servo with default values
func DefaultServo() Servo {
	return Servo{
		uuid.New().String(),
		"",
		-1,
		servoAvg,
		servoHardMin,
		servoHardMax,
		servoHardMin,
		servoHardMax,
		false,
	}
}

// SetDefaultPosition sets servo defaultPosition
func (servo *Servo) SetDefaultPosition(defaultPosition int) error {
	if defaultPosition < servo.Min {
		return errors.New("Servo.SetDefaultPosition() error: given position is lower than min")
	}
	if defaultPosition > servo.Max {
		return errors.New("Servo.SetDefaultPosition() error: given position is higher than max")
	}
	servo.DefaultPosition = defaultPosition
	return nil
}

// SetMin sets minimum position for servo
func (servo *Servo) SetMin(min int) error {
	if min < servo.HardMin {
		return errors.New("Servo.SetMin() error: given min is lower than hardMin")
	}
	servo.Min = min
	return nil
}

// SetMax sets maximum position for servo
func (servo *Servo) SetMax(max int) error {
	if max > servo.HardMax {
		return errors.New("Servo.SetMax() error: given max is higher than hardMax")
	}
	servo.Max = max
	return nil
}
