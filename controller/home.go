package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HomeHandler serves the home page
func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	writeJSONResponse(w, http.StatusOK, "Puppet Stage Server")
}
