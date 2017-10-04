package files

import (
	"encoding/json"
	"io/ioutil"
)

// Configuration type holds configuration values of the server
type Configuration struct {
	PuppetsDirectory string `json:"puppetsDirectory"`
	StagesDirectory  string `json:"stagesDirectory"`
}

// Conf retains server configuration loaded from configuration file
var Conf = loadConf()

func loadConf() Configuration {
	bytes, err := ioutil.ReadFile("./puppet-stage.conf.json")
	if err != nil {
		panic(err)
	}

	var conf Configuration
	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		panic(err)
	}

	return conf
}
