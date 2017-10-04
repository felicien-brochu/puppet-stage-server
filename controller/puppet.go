package controller

import (
	"felicien/puppet-server/files"
	"felicien/puppet-server/model"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetPuppetHandler sends a representation of the current puppet
func GetPuppetHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if isJSONRequest(r) {
		getPuppetJSONHandler(w, r, params)
	} else {
		getPuppetHTMLHandler(w, r, params)
	}
}

func getPuppetJSONHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	puppet := model.GetCurrentPuppet()
	if puppet == nil {
		writeJSONError(w, http.StatusNotFound, "No current puppet")
	} else {
		writeJSONResponse(w, http.StatusOK, *puppet)
	}
}

func getPuppetHTMLHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	http.ServeFile(w, r, "./html/puppet.html")
}

// ListPuppetsHandler lists puppets saved on server
func ListPuppetsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	puppets, err := files.ListPuppets()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	} else {
		writeJSONResponse(w, http.StatusOK, puppets)
	}
}

// CreatePuppetHandler creates a new puppet
func CreatePuppetHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	if name == "" {
		writeJSONResponse(w, http.StatusBadRequest, "No name")
	}

	// TODO check name conflicts with existing puppets on Server
	// TODO check if there is no current Puppet

	puppet := model.CreatePuppet(name)
	log.Printf("Puppet created: %v\n", *puppet)
	writeJSONResponse(w, http.StatusCreated, *puppet)
}
