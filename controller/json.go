package controller

import (
	"encoding/json"
	"net/http"
)

type jsonError struct {
	status  int
	message string
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	messageStruct := jsonError{status, message}
	writeJSONResponse(w, status, messageStruct)
}

func writeJSONResponse(w http.ResponseWriter, status int, object interface{}) {
	body, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	w.Write(body)
	// fmt.Fprintln(w, string(body))
}
