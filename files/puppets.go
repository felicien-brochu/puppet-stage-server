package files

import (
	"encoding/json"
	"felicien/puppet-server/model"
	"io/ioutil"
	"log"
)

// ListPuppets retrieves puppets from puppetsDirectory
func ListPuppets() ([]model.Puppet, error) {
	files, err := ioutil.ReadDir(Conf.PuppetsDirectory)
	if err != nil {
		return nil, err
	}

	var puppet *model.Puppet
	var puppets = make([]model.Puppet, 0)
	for _, file := range files {
		puppet, err = getPuppetFromFile(file.Name())
		if err != nil {
			log.Printf("ListPuppets() error: %s\n", err)
			continue
		}
		puppets = append(puppets, *puppet)
	}

	return puppets, nil
}

// GetPuppet returns the puppet with the given name if present in puppetsDirectory
func GetPuppet(name string) (*model.Puppet, error) {
	fileName := name + ".puppet.json"
	return getPuppetFromFile(fileName)
}

func getPuppetFromFile(fileName string) (*model.Puppet, error) {
	path := Conf.PuppetsDirectory + "/" + fileName
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var puppet model.Puppet
	err = json.Unmarshal(raw, &puppet)
	if err != nil {
		return nil, err
	}

	return &puppet, nil
}
