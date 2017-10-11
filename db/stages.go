package db

import (
	"encoding/json"
	"felicien/puppet-server/model"

	"github.com/garyburd/redigo/redis"
)

// GetStage returns the stage with the given id if present in db
func GetStage(id string) (*model.Stage, error) {
	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", "stage:"+id))
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	stageJSON, err := redis.Bytes(conn.Do("GET", "stage:"+id))
	if err != nil {
		return nil, err
	}

	var stage model.Stage
	err = json.Unmarshal(stageJSON, &stage)
	if err != nil {
		return nil, err
	}

	return &stage, nil
}

// ListStages retrieves stages from redis
func ListStages() ([]model.Stage, error) {
	conn := pool.Get()
	defer conn.Close()

	stageKeys, err := redis.Values(conn.Do("KEYS", "stage:*"))
	if err != nil {
		return nil, err
	}
	if len(stageKeys) == 0 {
		return make([]model.Stage, 0), nil
	}
	stagesJSON, err := redis.ByteSlices(conn.Do("MGET", stageKeys...))
	if err != nil {
		return nil, err
	}

	var stages []model.Stage

	for _, stageJSON := range stagesJSON {
		var stage model.Stage
		err = json.Unmarshal(stageJSON, &stage)
		if err != nil {
			return nil, err
		}

		stages = append(stages, stage)
	}

	return stages, nil
}

// SaveStage saves a stage (CREATE or UPDATE)
func SaveStage(stage model.Stage) error {
	stageJSON, err := json.Marshal(stage)
	if err != nil {
		return err
	}

	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", "stage:"+stage.ID, stageJSON)
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

	_, err = conn.Do("DEL", "stage:"+stageID)
	if err != nil {
		return nil, err
	}
	return stage, nil
}
