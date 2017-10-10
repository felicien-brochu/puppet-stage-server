package model

import (
	"github.com/google/uuid"
)

// Puppet represents a puppet configuration
type Puppet struct {
	ID     string           `json:"id"`
	Name   string           `json:"name"`
	Boards map[string]Board `json:"boards"`
}

// CreatePuppet creates a new current Puppet stores it and returns it
func CreatePuppet(name string) Puppet {
	return Puppet{
		uuid.New().String(),
		name,
		make(map[string]Board),
	}
}
