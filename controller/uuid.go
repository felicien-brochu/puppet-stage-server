package controller

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// GetUUIDsHandler returns a defined number of UUIDs
func GetUUIDsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	number, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "n parameter is not a valid positive integer")
	}

	uuids := make([]string, 0)
	for i := 0; i < number; i++ {
		uuids = append(uuids, uuid.New().String())
	}
	writeJSONResponse(w, http.StatusOK, uuids)
}
