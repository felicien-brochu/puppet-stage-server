package puppet

// Puppet represents a puppet configuration
type Puppet struct {
	Name   string
	Boards []*Board
}

// CreatePuppet creates a new current Puppet stores it and returns it
func CreatePuppet(name string) Puppet {
	puppet := new(Puppet)
	puppet.Name = name
	store.Puppet = puppet

	return *puppet
}
