package controller

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
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
	buf := make([]byte, 1024)
	runtime.Stack(buf, false)
	log.Printf("Recover from panic: %v\n%v\n", err, string(buf))
	writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %v", err))
}
