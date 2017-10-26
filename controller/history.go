package controller

import (
	"encoding/json"
	"felicien/puppet-server/db"
	"felicien/puppet-server/model"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// GetStageHistoryHandler returns revisions of a stage (stageID) from a particular revision (if from is "" activeRevision
// will be used instead) taking prev revisions before and next revisions after.
func GetStageHistoryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	stageID := params.ByName("stageID")
	if stageID == "" {
		writeJSONError(w, http.StatusNotFound, "No stage ID in request")
		return
	}

	prev, err := strconv.Atoi(r.URL.Query().Get("prev"))
	if err != nil || prev < 0 {
		writeJSONError(w, http.StatusBadRequest, "prev parameter is not a valid positive integer")
		return
	}

	next, err := strconv.Atoi(r.URL.Query().Get("next"))
	if err != nil || prev < 0 {
		writeJSONError(w, http.StatusBadRequest, "next parameter is not a valid positive integer")
		return
	}

	from := r.URL.Query().Get("from")

	stageHistory, err := db.GetStageHistory(stageID, from, prev, next)
	if err != nil {
		panic(err)
	}
	if stageHistory == nil {
		writeJSONError(w, http.StatusNotFound, "Stage not found or from parameter does not correspond to an existing revision")
		return
	}

	writeJSONResponse(w, http.StatusOK, stageHistory)
}

// SaveStageHistoryHandler saves a partial stage history
func SaveStageHistoryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	stageID := params.ByName("stageID")
	if stageID == "" {
		writeJSONError(w, http.StatusNotFound, "No stage ID in request")
		return
	}

	var body saveHistoryBody
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Request body JSON not well formatted: %v", err))
		return
	}

	err = db.UpdateStageHistory(stageID, body.StartRevisionID, body.ActiveRevisionID, body.Revisions)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, "OK")
}

type saveHistoryBody struct {
	StartRevisionID  string                `json:"startRevisionID"`
	ActiveRevisionID string                `json:"activeRevisionID"`
	Revisions        []model.StageRevision `json:"revisions"`
}

// SaveHistoryActiveRevisionHandler saves thes history active revision
func SaveHistoryActiveRevisionHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	stageID := params.ByName("stageID")
	if stageID == "" {
		writeJSONError(w, http.StatusNotFound, "No stage ID in request")
		return
	}

	var body saveActiveRevisionBody
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Request body JSON not well formatted: %v", err))
		return
	}

	err = db.UpdateStageHistoryActiveRevision(stageID, body.ActiveRevisionID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, "OK")
}

type saveActiveRevisionBody struct {
	ActiveRevisionID string `json:"activeRevisionID"`
}
