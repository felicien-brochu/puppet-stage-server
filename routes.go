package main

import (
	"felicien/puppet-server/controller"

	"github.com/julienschmidt/httprouter"
)

func getRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/websocket", controller.WebsocketHandler)

	router.NotFound = controller.NotFoundHandler{}
	router.MethodNotAllowed = controller.MethodNotAllowedHandler{}

	router.GET("/", controller.HomeHandler)
	router.GET("/puppet", controller.GetPuppetHandler)
	router.GET("/puppets", controller.ListPuppetsHandler)
	router.PUT("/puppet/:name", controller.CreatePuppetHandler)

	return router
}
