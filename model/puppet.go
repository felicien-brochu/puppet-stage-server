package model

import (
	"github.com/google/uuid"
)

// Puppet represents a puppet configuration
type Puppet struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Boards []*Board `json:"boards"`
}

// CreatePuppet creates a new current Puppet stores it and returns it
func CreatePuppet(name string) Puppet {
	var puppet Puppet
	puppet.ID = uuid.New().String()
	puppet.Name = name
	puppet.Boards = make([]*Board, 0)

	return puppet
}
