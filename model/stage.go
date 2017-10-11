package model

import (
	"github.com/google/uuid"
)

// Stage is the name of a project in Puppet Stage
type Stage struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Puppet    *Puppet    `json:"puppet"`
	Sequences []Sequence `json:"sequences"`
	Duration  Duration   `json:"duration"`
	History   []Stage    `json:"history"`
}

// Sequence defines a sequence of values over time
type Sequence interface {
	ValueAt(t Time) (float64, error)
	StartTime() Time
	Duration() Duration
}

// DriverSequence is a sequence that can drive a servo
type DriverSequence interface {
	Sequence
	Servo() *Servo
}

// NewStage returns a new stage
func NewStage(name string) Stage {
	return Stage{
		ID:        uuid.New().String(),
		Name:      name,
		Sequences: make([]Sequence, 0),
		Duration:  10 * Second,
		History:   make([]Stage, 0),
	}
}
