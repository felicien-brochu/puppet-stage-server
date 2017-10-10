package controller

import (
	"encoding/json"
	"felicien/puppet-server/db"
	"felicien/puppet-server/model"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetPuppetHandler sends a representation of the current puppet
func GetPuppetHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("puppetID")
	if id == "" {
		writeJSONError(w, http.StatusNotFound, "No ID in request")
		return
	}

	puppet, err := db.GetPuppet(id)
	if err != nil {
		panic(err)
	}
	if puppet == nil {
		writeJSONError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	writeJSONResponse(w, http.StatusOK, puppet)
}

// ListPuppetsHandler lists puppets saved on server
func ListPuppetsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	puppets, err := db.ListPuppets()
	if err != nil {
		panic(err)
	} else {
		writeJSONResponse(w, http.StatusOK, puppets)
	}
}

// CreatePuppetHandler creates a new puppet
func CreatePuppetHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var puppet model.Puppet
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &puppet)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Puppet JSON not well formatted: %v", err))
		return
	}
	if puppet.Name == "" {
		writeJSONError(w, http.StatusBadRequest, "Puppet JSON must contain a name")
	}
	puppet = model.CreatePuppet(puppet.Name)
	err = db.SavePuppet(puppet)
	if err != nil {
		panic(err)
	}

	log.Printf("Puppet created: %v\n", puppet)
	writeJSONResponse(w, http.StatusCreated, puppet)
}

// UpdatePuppetHandler creates a new puppet
func UpdatePuppetHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var puppet model.Puppet
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &puppet)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Puppet JSON not well formatted: %v", err))
		return
	}
	err = db.SavePuppet(puppet)
	if err != nil {
		panic(err)
	}

	writeJSONResponse(w, http.StatusOK, puppet)
}
