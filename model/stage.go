package model

import (
	"math"

	"github.com/google/uuid"
)

// Stage is the name of a project in Puppet Stage
type Stage struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	PuppetID  string           `json:"puppetID"`
	Sequences []DriverSequence `json:"sequences"`
	Duration  Duration         `json:"duration"`
}

// Sequence defines a sequence of values over time
type Sequence interface {
	GetID() string
	StartTime() Time
	TotalDuration() Duration
	ValueAt(t Time) (float64, error)
}

// DriverSequence is a sequence that can drive a servo.
// First level sequence, it determines its values by its
// subsequences.
type DriverSequence struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	ServoID     string          `json:"servoID"`
	Expanded    bool            `json:"expanded"`
	Color       int             `json:"color"`
	PlayEnabled bool            `json:"playEnabled"`
	Sequences   []BasicSequence `json:"sequences"`
}

// InitStage inits a new stage
func InitStage(stage Stage) Stage {
	stage.ID = uuid.New().String()
	stage.Sequences = make([]DriverSequence, 0)
	stage.Duration = 10 * Second

	return stage
}

// GetFrameAt returns the values of all sequences at a time t.
func (stage *Stage) GetFrameAt(t Time, preview bool) map[string]float64 {
	var frame = make(map[string]float64)
	for _, driverSequence := range stage.Sequences {
		frame[driverSequence.ServoID] = driverSequence.GetValueAt(t, preview)
	}
	return frame
}

// GetValueAt returns the value of the driverSequence at time t.
func (driverSequence *DriverSequence) GetValueAt(t Time, preview bool) float64 {
	var value = math.NaN()
	if !preview && !driverSequence.PlayEnabled {
		return value
	}

	for _, basicSequence := range driverSequence.Sequences {
		if preview && !basicSequence.PreviewEnabled || !preview && !basicSequence.PlayEnabled {
			continue
		}

		if basicSequence.Start.Before(t) && (basicSequence.Start + Time(basicSequence.Duration)).After(t) {
			value = basicSequence.ValueAt(t)
			break
		}
	}
	return value
}
