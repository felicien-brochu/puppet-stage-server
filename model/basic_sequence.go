package model

import (
	"fmt"
)

// BasicSequence is a monovalued sequence defined by Bezier curves
type BasicSequence struct {
	start    Time
	duration Duration
	Curves   []BezierCurve
}

const (
	// BezierTimePrecision time precision of Bezier curves. Because values
	// are approximated (faster) this precision is the upper bound of the approximation.
	BezierTimePrecision Duration = Nanosecond
)

// Point represents a time point
type Point struct {
	T Time
	V float64
}

// BezierCurve defines a cubic bezier curve
type BezierCurve struct {
	P1 Point
	C1 Point
	P2 Point
	C2 Point
}

// NewBasicSequence returns a new basic sequence
func NewBasicSequence() BasicSequence {
	return BasicSequence{
		0,
		10 * Second,
		make([]BezierCurve, 0),
	}
}

// NewBezierCurve returns a new BezierCurve with control points on the value points
func NewBezierCurve(a Point, b Point) BezierCurve {
	return BezierCurve{
		a,
		a,
		b,
		b,
	}
}

// StartTime returns the starting time of the sequence
func (sequence *BasicSequence) StartTime() Time {
	return sequence.start
}

// Duration returns the duration of the sequence
func (sequence *BasicSequence) Duration() Duration {
	return sequence.duration
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
	for _, curve := range sequence.Curves {
		if curve.P1.T.Before(t) && curve.P2.T.After(t) {
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
