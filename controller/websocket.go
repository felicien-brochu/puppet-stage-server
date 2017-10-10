package controller

import (
	"encoding/json"
	"felicien/puppet-server/db"
	"felicien/puppet-server/model"
	"felicien/puppet-server/players"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

// ServoPositionWebsocketHandler handles websocket for testing servo on a given puppet
func ServoPositionWebsocketHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	puppetID := params.ByName("puppetID")
	puppet, err := db.GetPuppet(puppetID)
	if err != nil {
		panic(err)
	}
	if puppet == nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("No puppet for id '%s'.", puppetID))
		return
	}

	player := players.GetPuppetPlayer(*puppet)
	if player == nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("No puppet player for id '%s'.", puppetID))
		return
	}

	boardID := params.ByName("boardID")
	boardPlayer := player.GetBoardPlayer(boardID)
	if boardPlayer == nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("No board player for id '%s'.", boardID))
		return
	}

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	defer connection.Close()
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			panic(err)
		}

		var jsonCommand jsonPositionCommand
		err = json.Unmarshal(message, &jsonCommand)
		if err != nil {
			continue
		}
		var positionCommand = model.PositionCommand{
			Addr:     jsonCommand.Addr,
			Position: jsonCommand.Position,
		}

		boardPlayer.AddCommand(positionCommand)
	}
}

type jsonPositionCommand struct {
	Addr     int
	Position int
}

// // WebsocketHandler handles websocket messages
// func WebsocketHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	connection, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Print("upgrade:", err)
// 		return
// 	}
//
// 	defer connection.Close()
// 	for {
// 		messageType, message, err := connection.ReadMessage()
// 		if err != nil {
// 			log.Println("read:", err)
// 			return
// 		}
//
// 		err = connection.WriteMessage(messageType, message)
// 		if err != nil {
// 			log.Println("write:", err)
// 			return
// 		}
// 	}
// }
