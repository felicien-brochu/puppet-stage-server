package controller

import (
	"encoding/json"
	"felicien/puppet-server/db"
	"felicien/puppet-server/drivers"
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

	driver := drivers.GetPuppetDriver(*puppet)
	if driver == nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("No puppet driver for id '%s'.", puppetID))
		return
	}

	boardID := params.ByName("boardID")
	boardDriver := driver.GetBoardDriver(boardID)
	if boardDriver == nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("No board driver for id '%s'.", boardID))
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
		var positionCommand = drivers.PositionCommand{
			Addr:     jsonCommand.Addr,
			Position: jsonCommand.Position,
		}

		boardDriver.AddCommand(positionCommand)
	}
}

type jsonPositionCommand struct {
	Addr     int
	Position int
}

// PuppetPlayerWebsocketHandler handles the puppet player commands
func PuppetPlayerWebsocketHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	puppetID := params.ByName("puppetID")
	puppet, err := db.GetPuppet(puppetID)
	if err != nil {
		panic(err)
	}
	if puppet == nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("No puppet for id '%s'.", puppetID))
		return
	}

	driver := drivers.GetPuppetDriver(*puppet)
	if driver == nil {
		driver, err = drivers.AddPuppetDriver(*puppet)
		if err != nil {
			panic(err)
		}
	}

	err = driver.Start()
	if err != nil {
		panic(err)
	}
	defer driver.Stop()

	var player = players.NewPuppetPlayer(*puppet, driver)

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

		var jsonCommand jsonPlayerCommand
		err = json.Unmarshal(message, &jsonCommand)
		if err != nil {
			continue
		}

		if jsonCommand.Type == "preview" {
			var previewCommand jsonPreviewCommand
			err = json.Unmarshal(jsonCommand.Body, &previewCommand)
			if err != nil {
				continue
			}

			player.PreviewServo(previewCommand.ServoID, previewCommand.Value)
		} else if jsonCommand.Type == "play" {
			var playCommand jsonPlayCommand
			err = json.Unmarshal(jsonCommand.Body, &playCommand)
			if err != nil {
				continue
			}

			go playStage(player, playCommand.Stage, playCommand.PlayStart, connection)
		}
	}
}

func playStage(player players.PuppetPlayer, stage model.Stage, playStart model.Time, connection *websocket.Conn) {
	var playerState = make(chan string)
	player.PlayStage(stage, playStart, playerState)

	for {
		state, more := <-playerState
		connection.WriteMessage(websocket.TextMessage, []byte(state))
		if !more {
			return
		}
	}
}

type jsonPlayerCommand struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

type jsonPreviewCommand struct {
	ServoID string  `json:"servoID"`
	Value   float64 `json:"value"`
}

type jsonPlayCommand struct {
	Stage     model.Stage `json:"stage"`
	PlayStart model.Time  `json:"playStart"`
}
