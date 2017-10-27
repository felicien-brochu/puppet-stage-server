package model

import (
	"fmt"

	"github.com/google/uuid"
)

// BasicSequence is a monovalued sequence defined by Bezier curves
type BasicSequence struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Start        Time       `json:"start"`
	Duration     Duration   `json:"duration"`
	DefaultValue float64    `json:"defaultValue"`
	Keyframes    []Keyframe `json:"keyframes"`
	Slave        bool       `json:"slave"`
}

const (
	// BezierTimePrecision time precision of Bezier curves. Because values
	// are approximated (faster) this precision is the upper bound of the approximation.
	BezierTimePrecision Duration = Nanosecond
)

// Keyframe defines a keyframe that can be used to construct cubic bezier curves
type Keyframe struct {
	P  Point `json:"p"`
	C1 Point `json:"c1"`
	C2 Point `json:"c2"`
}

// BezierCurve defines a cubic bezier curve
type BezierCurve struct {
	P1 Point `json:"p1"`
	C1 Point `json:"c1"`
	P2 Point `json:"p2"`
	C2 Point `json:"c2"`
}

// Point represents a time point
type Point struct {
	T Time    `json:"t"`
	V float64 `json:"v"`
}

// NewBasicSequence returns a new basic sequence
func NewBasicSequence() BasicSequence {
	return BasicSequence{
		uuid.New().String(),
		"",
		0,
		10 * Second,
		0,
		make([]Keyframe, 0),
		false,
	}
}

// NewBezierCurve returns a new BezierCurve
func NewBezierCurve(p1 Point, c1 Point, p2 Point, c2 Point) BezierCurve {
	return BezierCurve{
		p1,
		c1,
		p2,
		c2,
	}
}

// KeyframesToBezierCurve returns a new BezierCurve from 2 Keyframes
func KeyframesToBezierCurve(a Keyframe, b Keyframe) BezierCurve {
	return BezierCurve{
		a.P,
		a.C2,
		b.P,
		b.C1,
	}
}

// GetID returns the sequence id.
func (sequence *BasicSequence) GetID() string {
	return sequence.ID
}

// StartTime returns the starting time of the sequence
func (sequence *BasicSequence) StartTime() Time {
	return sequence.Start
}

// TotalDuration returns the duration of the sequence
func (sequence *BasicSequence) TotalDuration() Duration {
	return sequence.Duration
}

// ValueAt returns the value of the sequence at the given time
func (sequence *BasicSequence) ValueAt(t Time) (float64, error) {
	curve := sequence.curveAt(t)
	if curve == nil {
		return -1, fmt.Errorf("No value at time %d", t)
	}

	return curve.ValueAt(t), nil
}

func (sequence *BasicSequence) curveAt(t Time) *BezierCurve {
	for i := 0; i < len(sequence.Keyframes)-1; i++ {
		keyframe1 := sequence.Keyframes[i]
		keyframe2 := sequence.Keyframes[i+1]

		if keyframe1.P.T.Before(t) && keyframe2.P.T.After(t) {
			curve := KeyframesToBezierCurve(keyframe1, keyframe2)
			return &curve
		}
	}
	return nil
}

// ValueAt returns the value of the curve at t with a precision of BezierTimePrecision
func (curve *BezierCurve) ValueAt(t Time) float64 {
	var progress = float64(0.5)
	var min = float64(0)
	var max = float64(1)
	var point Point

	for {
		point = curve.progressPointAt(progress)
		if point.T.AbsDiff(t) < BezierTimePrecision {
			break
		} else if point.T.Before(t) {
			min = progress
		} else {
			max = progress
		}
		progress = (max + min) / 2
	}
	return point.V
}

func (curve *BezierCurve) progressPointAt(progress float64) Point {
	var a = progressPoint(progress, curve.P1, curve.C1)
	var b = progressPoint(progress, curve.C1, curve.C2)
	var c = progressPoint(progress, curve.C2, curve.P2)
	var d = progressPoint(progress, a, b)
	var e = progressPoint(progress, b, c)
	return progressPoint(progress, d, e)
}

func progressPoint(progress float64, a Point, b Point) Point {
	return Point{
		T: Time(float64((b.T-a.T))*progress) + a.T,
		V: (b.V-a.V)*progress + a.V,
	}
}
