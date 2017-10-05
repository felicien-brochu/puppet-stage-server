package main

import (
	"felicien/puppet-server/controller"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func getRouter() http.Handler {
	router := httprouter.New()

	router.GET("/websocket", controller.WebsocketHandler)

	router.NotFound = controller.NotFoundHandler{}
	router.MethodNotAllowed = controller.MethodNotAllowedHandler{}
	router.PanicHandler = controller.PanicHandler

	router.GET("/", controller.HomeHandler)
	router.GET("/puppet/:id", controller.GetPuppetHandler)
	router.GET("/puppets", controller.ListPuppetsHandler)
	router.PUT("/puppet", controller.CreatePuppetHandler)
	router.PUT("/puppet/:id?", controller.UpdatePuppetHandler)

	handler := cors.AllowAll().Handler(router)

	return handler
}
