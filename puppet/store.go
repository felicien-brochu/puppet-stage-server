package puppet

// Storage stores runtime objects of the server
type Storage struct {
	Puppet *Puppet
}

// store stores the model of the server
var store = Storage{}
