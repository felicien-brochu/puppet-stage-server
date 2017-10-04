package model

// Puppet represents a puppet configuration
type Puppet struct {
	Name   string   `json:"name"`
	Boards []*Board `json:"boards"`
}

// CreatePuppet creates a new current Puppet stores it and returns it
func CreatePuppet(name string) *Puppet {
	puppet := new(Puppet)
	puppet.Name = name
	store.Puppet = puppet

	return puppet
}

// GetCurrentPuppet returns the current puppet
func GetCurrentPuppet() *Puppet {
	return store.Puppet
}
