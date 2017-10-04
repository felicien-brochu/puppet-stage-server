package controller

import (
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
