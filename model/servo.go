package model

import "errors"

// Servo describes a servomotor
type Servo struct {
	Name            string
	Addr            int
	DefaultPosition int
	HardMin         int
	HardMax         int
	Min             int
	Max             int
}

const (
	servoHardMin = 130
	servoHardMax = 470
	servoAvg     = (servoHardMin + servoHardMax) / 2
)

// NewServo constructs a new Servo with default values
func NewServo() Servo {
	return Servo{
		"",
		-1,
		servoAvg,
		servoHardMin,
		servoHardMax,
		servoHardMin,
		servoHardMax,
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
