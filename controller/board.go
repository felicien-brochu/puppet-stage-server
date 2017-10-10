package controller

import (
	"felicien/puppet-server/db"
	"felicien/puppet-server/model"
	"felicien/puppet-server/players"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// NewBoardHandler returns default new board with new unique ID
func NewBoardHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var board = model.DefaultBoard()
	writeJSONResponse(w, http.StatusOK, board)
}

// StartBoardsHandler starts boards connections' if not done already
func StartBoardsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
		var err error
		player, err = players.AddPuppetPlayer(*puppet)
		if err != nil {
			panic(err)
		}
	}

	player.Start()
	writeJSONResponse(w, http.StatusOK, "OK")
}

// StartBoardHandler starts board connection if not done already
func StartBoardHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	puppetID := params.ByName("puppetID")
	puppet, err := db.GetPuppet(puppetID)
	if err != nil {
		panic(err)
	}
	if puppet == nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("No puppet for id '%s'.", puppetID))
		return
	}
	boardID := params.ByName("boardID")
	board, ok := puppet.Boards[boardID]
	if !ok {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("No board for id '%s'.", boardID))
		return
	}

	player := players.GetPuppetPlayer(*puppet)
	if player == nil {
		player, err = players.AddPuppetPlayer(*puppet)
		if err != nil {
			panic(err)
		}
	}

	err = player.StartBoard(board.ID)
	if err != nil {
		panic(err)
	}
	writeJSONResponse(w, http.StatusOK, "OK")
}
