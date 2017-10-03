package controller

import "net/http"

func isHTMLRequest(request *http.Request) bool {
	accept := request.Header.Get("Accept")
	return accept == "" || strings.Contains(accept, "*/*") || strings.Contains(accept, "text/html")
}

func isJSONRequest(request *http.Request) bool {
	accept := request.Header.Get("Accept")
	return accept == "" || strings.Contains(accept, "*/*") || strings.Contains(accept, "application/json")
}
