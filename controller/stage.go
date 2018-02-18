package controller

import (
	"encoding/json"
	"felicien/puppet-server/db"
	"felicien/puppet-server/model"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// GetStageHandler sends a representation of the current stage
func GetStageHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	stageID := params.ByName("stageID")
	if stageID == "" {
		writeJSONError(w, http.StatusNotFound, "No ID in request")
		return
	}

	stage, err := db.GetStage(stageID)
	if err != nil {
		panic(err)
	}
	if stage == nil {
		writeJSONError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	writeJSONResponse(w, http.StatusOK, stage)
}

// ListStagesHandler lists stages saved on server
func ListStagesHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	stages, err := db.ListStages()
	if err != nil {
		panic(err)
	} else {
		writeJSONResponse(w, http.StatusOK, stages)
	}
}

// CreateStageHandler creates a new stage
func CreateStageHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var stage model.Stage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &stage)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Stage JSON not well formatted: %v", err))
		return
	}
	if stage.Name == "" {
		writeJSONError(w, http.StatusBadRequest, "Stage JSON must contain a name")
		return
	}
	if stage.PuppetID == "" {
		writeJSONError(w, http.StatusBadRequest, "Stage JSON must contain a puppetID")
		return
	}
	stage = InitStage(stage)
	err = db.CreateStage(stage)
	if err != nil {
		panic(err)
	}

	log.Printf("Stage created: %v\n", stage)
	writeJSONResponse(w, http.StatusCreated, stage)
}

// InitStage inits a new stage
func InitStage(stage model.Stage) model.Stage {
	stage.ID = uuid.New().String()
	stage.Sequences = make([]model.DriverSequence, 0)
	stage.Duration = 10 * model.Second

	// Create sequences for each servo
	if stage.PuppetID != "" {
		puppet, err := db.GetPuppet(stage.PuppetID)
		if err != nil {
			panic(err)
		}

		if puppet != nil {
			var color = 0
			for _, board := range puppet.Boards {
				for _, servo := range board.Servos {
					var defaultValue float64
					if servo.Inverted {
						defaultValue = float64(servo.Max-servo.DefaultPosition) / float64(servo.Max-servo.Min) * 100.
					} else {
						defaultValue = float64(servo.DefaultPosition-servo.Min) / float64(servo.Max-servo.Min) * 100.
					}

					var basicSequence = model.BasicSequence{
						ID:             uuid.New().String(),
						Name:           "Main",
						Start:          0,
						Duration:       stage.Duration,
						DefaultValue:   defaultValue,
						Slave:          false,
						PlayEnabled:    true,
						PreviewEnabled: false,
						ShowGraph:      false,
						Keyframes:      make([]model.Keyframe, 0),
					}
					var driverSequence = model.DriverSequence{
						ID:          uuid.New().String(),
						Name:        servo.Name,
						ServoID:     servo.ID,
						Expanded:    true,
						Color:       color,
						PlayEnabled: true,
						Sequences:   []model.BasicSequence{basicSequence},
					}
					color++
					stage.Sequences = append(stage.Sequences, driverSequence)
				}
			}
		}
	}

	return stage
}

// DeleteStageHandler deletes a stage
func DeleteStageHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	stageID := params.ByName("stageID")
	stage, err := db.DeleteStage(stageID)
	if err != nil {
		panic(err)
	}
	if stage == nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("No stage for id '%s'", stageID))
		return
	}

	writeJSONResponse(w, http.StatusOK, stage)
}
