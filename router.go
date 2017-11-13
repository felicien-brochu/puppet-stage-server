package main

import (
	"felicien/puppet-server/controller"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func getRouter() http.Handler {
	router := httprouter.New()

	router.GET("/puppet/:puppetID/board/:boardID/websocket", controller.ServoPositionWebsocketHandler)
	router.GET("/puppet/:puppetID/player", controller.PuppetPlayerWebsocketHandler)

	router.NotFound = controller.NotFoundHandler{}
	router.MethodNotAllowed = controller.MethodNotAllowedHandler{}
	router.PanicHandler = controller.PanicHandler

	router.GET("/", controller.HomeHandler)
	router.GET("/uuids", controller.GetUUIDsHandler)

	router.GET("/board/new", controller.NewBoardHandler)
	router.GET("/servo/new", controller.NewServoHandler)

	router.GET("/puppet/:puppetID", controller.GetPuppetHandler)
	router.GET("/puppets", controller.ListPuppetsHandler)
	router.PUT("/puppet", controller.CreatePuppetHandler)
	router.PUT("/puppet/:puppetID?", controller.UpdatePuppetHandler)
	router.DELETE("/puppet/:puppetID", controller.DeletePuppetHandler)

	router.POST("/puppet/:puppetID/boards/start", controller.StartBoardsHandler)
	router.POST("/puppet/:puppetID/board/:boardID/start", controller.StartBoardHandler)

	router.GET("/stage/:stageID", controller.GetStageHandler)
	router.GET("/stages", controller.ListStagesHandler)
	router.PUT("/stage", controller.CreateStageHandler)
	router.DELETE("/stage/:stageID", controller.DeleteStageHandler)
	router.GET("/stage/:stageID/history", controller.GetStageHistoryHandler)
	router.PUT("/stage/:stageID/history", controller.SaveStageHistoryHandler)
	router.PUT("/stage/:stageID/history/activeRevision", controller.SaveHistoryActiveRevisionHandler)

	router.GET("/audio/:fileName", controller.GetAudioFile)
	router.PUT("/audio/:fileName", controller.HandleAudioFileUpload)

	handler := cors.AllowAll().Handler(router)

	return handler
}
