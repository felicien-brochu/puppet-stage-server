package controller

import (
	"felicien/puppet-server/model"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// NewServoHandler returns default new servo with new unique ID
func NewServoHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var servo = model.DefaultServo()
	writeJSONResponse(w, http.StatusOK, servo)
}
