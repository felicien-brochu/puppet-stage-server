package controller

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}} // use default options

// WebsocketHandler handles websocket messages
func WebsocketHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer connection.Close()
	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		err = connection.WriteMessage(messageType, message)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
