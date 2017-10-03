package controller

import (
	"felicien/puppet-server/puppet"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GetPuppetHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}

// CreatePuppetHandler creates a new puppet
func CreatePuppetHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	if name == "" {
		writeJSONResponse(w, http.StatusBadRequest, "No name")
	}

	// TODO check name conflicts with existing puppets on Server
	// TODO check if there is no current Puppet

	puppet := puppet.CreatePuppet(name)
	log.Printf("Puppet created: %v\n", puppet)
	writeJSONResponse(w, http.StatusCreated, puppet)
}
