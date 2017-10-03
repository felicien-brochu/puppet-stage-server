package main

import (
	"felicien/puppet-server/controller"

	"github.com/julienschmidt/httprouter"
)

func getRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/websocket", controller.WebsocketHandler)

	router.GET("/", controller.HomeHandler)
	router.PUT("/puppet/:name", controller.CreatePuppetHandler)

	return router
}
