package db

import (
	"encoding/json"
	"felicien/puppet-server/model"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/google/uuid"
)

// GetStage returns the stage with the given id if present in db
func GetStage(stageID string) (*model.Stage, error) {
	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", fmt.Sprintf("stageHistory:%s", stageID)))
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	stageHistoryJSON, err := redis.Bytes(conn.Do("GET", fmt.Sprintf("stageHistory:%s", stageID)))
	if err != nil {
		return nil, err
	}

	var stageHistory model.StageHistoryRef
	err = json.Unmarshal(stageHistoryJSON, &stageHistory)
	if err != nil {
		return nil, err
	}

	stageRevisionJSON, err := redis.Bytes(conn.Do("GET", fmt.Sprintf("stageRevision:%s", stageHistory.ActiveRevision)))
	if err != nil {
		return nil, err
	}

	var stageRevision model.StageRevision
	err = json.Unmarshal(stageRevisionJSON, &stageRevision)
	if err != nil {
		return nil, err
	}

	return &stageRevision.Stage, nil
}

// ListStages retrieves active revision stages from redis
func ListStages() ([]model.Stage, error) {
	conn := pool.Get()
	defer conn.Close()

	stageHistoryKeys, err := redis.Values(conn.Do("KEYS", "stageHistory:*"))
	if err != nil {
		return nil, err
	}
	if len(stageHistoryKeys) == 0 {
		return make([]model.Stage, 0), nil
	}
	stageHistoriesJSON, err := redis.ByteSlices(conn.Do("MGET", stageHistoryKeys...))
	if err != nil {
		return nil, err
	}

	var activeRevisions []interface{}

	for _, stageHistoryJSON := range stageHistoriesJSON {
		var stageHistory model.StageHistoryRef
		err = json.Unmarshal(stageHistoryJSON, &stageHistory)
		if err != nil {
			return nil, err
		}
		activeRevisions = append(activeRevisions, fmt.Sprintf("stageRevision:%s", stageHistory.ActiveRevision))
	}

	revisionsJSON, err := redis.ByteSlices(conn.Do("MGET", activeRevisions...))
	if err != nil {
		return nil, err
	}

	var stages []model.Stage

	for _, revisionJSON := range revisionsJSON {
		var revision model.StageRevision
		err = json.Unmarshal(revisionJSON, &revision)
		if err != nil {
			return nil, err
		}
		stages = append(stages, revision.Stage)
	}

	return stages, nil
}

// CreateStage creates a stage with its history
func CreateStage(stage model.Stage) error {
	var revision = model.StageRevision{
		ID:    uuid.New().String(),
		Stage: stage,
		Date:  time.Now(),
	}

	var revisions = []string{revision.ID}
	var archives = make([]string, 0)
	var stageHistory = model.StageHistoryRef{
		StageID:        stage.ID,
		ActiveRevision: revision.ID,
		Revisions:      revisions,
		Archives:       archives,
	}

	historyJSON, err := json.Marshal(stageHistory)
	if err != nil {
		return err
	}

	revisionJSON, err := json.Marshal(revision)
	if err != nil {
		return err
	}

	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("MSET",
		fmt.Sprintf("stageHistory:%s", stage.ID), historyJSON,
		fmt.Sprintf("stageRevision:%s", revision.ID), revisionJSON)
	if err != nil {
		return err
	}
	return nil
}

// DeleteStage delete a stage
func DeleteStage(stageID string) (*model.Stage, error) {
	stage, err := GetStage(stageID)
	if err != nil {
		return nil, err
	}
	if stage == nil {
		return nil, nil
	}

	conn := pool.Get()
	defer conn.Close()

	stageHistoryJSON, err := redis.Bytes(conn.Do("GET", fmt.Sprintf("stageHistory:%s", stageID)))
	if err != nil {
		return nil, err
	}

	var stageHistory model.StageHistoryRef
	err = json.Unmarshal(stageHistoryJSON, &stageHistory)
	if err != nil {
		return nil, err
	}
	var delKeys = []interface{}{fmt.Sprintf("stageHistory:%s", stageID)}

	for _, revisionID := range stageHistory.Revisions {
		delKeys = append(delKeys, fmt.Sprintf("stageRevision:%s", revisionID))
	}
	for _, revisionID := range stageHistory.Archives {
		delKeys = append(delKeys, fmt.Sprintf("stageRevision:%s", revisionID))
	}

	_, err = conn.Do("DEL", delKeys...)
	if err != nil {
		return nil, err
	}
	return stage, nil
}
