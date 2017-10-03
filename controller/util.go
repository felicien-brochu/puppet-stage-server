package controller

import (
	"net/http"
	"strings"
)

func isHTMLRequest(request *http.Request) bool {
	return getPreferedFormat(request) == "text/html"
}

func isJSONRequest(request *http.Request) bool {
	return getPreferedFormat(request) == "application/json"
}

func getPreferedFormat(request *http.Request) string {
	accept := request.Header.Get("Accept")
	formats := strings.Split(accept, ",")
	choice := "text/html"
	for _, format := range formats {
		if strings.Contains(format, "text/html") {
			break
		} else if strings.Contains(format, "application/json") {
			choice = "application/json"
			break
		}
	}
	return choice
}
