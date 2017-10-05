package controller

import (
	"fmt"
	"net/http"
)

// NotFoundHandler sends a JSON error message
type NotFoundHandler struct{}

// ServeHTTP sends a JSON error message
func (h NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSONError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

// MethodNotAllowedHandler sends a JSON error message
type MethodNotAllowedHandler struct{}

// ServeHTTP sends a JSON error message
func (h MethodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSONError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

// PanicHandler handles panic in router. Sends a JSON object describing the internal server error (status 500).
func PanicHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %v", err))
}
