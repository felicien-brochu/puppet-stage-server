package db

import (
	"encoding/json"
	"felicien/puppet-server/model"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// GetStageHistory returns revisions of a stage (stageID) from a particular revision (if from is "" activeRevision
// will be used instead) taking prev revisions before and next revisions after.
func GetStageHistory(stageID, from string, prev, next int) (*model.StageHistory, error) {
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

	var stageHistoryRef model.StageHistoryRef
	err = json.Unmarshal(stageHistoryJSON, &stageHistoryRef)
	if err != nil {
		return nil, err
	}

	if from == "" {
		from = stageHistoryRef.ActiveRevision
	}

	fromIndex := -1
	revisionIDs := append(stageHistoryRef.Archives, stageHistoryRef.Revisions...)
	for i, revisionID := range revisionIDs {
		if revisionID == from {
			fromIndex = i
			break
		}
	}

	if fromIndex < 0 {
		return nil, nil
	}

	startIndex := fromIndex - prev
	if startIndex < 0 {
		startIndex = 0
	}

	endIndex := fromIndex + next + 1
	if endIndex > len(revisionIDs) {
		endIndex = len(revisionIDs)
	}

	var revisionKeys []interface{}

	for _, revisionID := range revisionIDs[startIndex:endIndex] {
		revisionKeys = append(revisionKeys, fmt.Sprintf("stageRevision:%s", revisionID))
	}

	stageRevisionsJSON, err := redis.ByteSlices(conn.Do("MGET", revisionKeys...))
	if err != nil {
		return nil, err
	}

	var revisions []model.StageRevision
	for _, revisionJSON := range stageRevisionsJSON {
		var revision model.StageRevision
		err = json.Unmarshal(revisionJSON, &revision)
		if err != nil {
			return nil, err
		}
		revisions = append(revisions, revision)
	}

	var history = model.StageHistory{
		StageID:        stageID,
		ActiveRevision: stageHistoryRef.ActiveRevision,
		Revisions:      revisions,
	}

	return &history, nil
}

// UpdateStageHistory updates a stage history with additional revisions. Adds these revisions immediatly after
// startRevisionID. All revisions that exists after this one will be erased. It updates also the active revision.
func UpdateStageHistory(stageID, startRevisionID, activeRevisionID string, revisions []model.StageRevision) error {
	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", fmt.Sprintf("stageHistory:%s", stageID)))
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("No stage history for id '%s'", stageID)
	}

	stageHistoryJSON, err := redis.Bytes(conn.Do("GET", fmt.Sprintf("stageHistory:%s", stageID)))
	if err != nil {
		return err
	}

	var stageHistoryRef model.StageHistoryRef
	err = json.Unmarshal(stageHistoryJSON, &stageHistoryRef)
	if err != nil {
		return err
	}

	startIndex := -1
	oldRevisions := append(stageHistoryRef.Archives, stageHistoryRef.Revisions...)

	for i, revisionID := range oldRevisions {
		if revisionID == startRevisionID {
			startIndex = i
			break
		}
	}

	if startIndex < 0 {
		return fmt.Errorf("StartRevisionID '%s' not found", startRevisionID)
	}

	revisionIDs := make([]string, 0)
	for _, revision := range revisions {
		revisionIDs = append(revisionIDs, revision.ID)
	}

	revisionsToDelete := oldRevisions[startIndex+1:]
	revisionsStartIndex := startIndex - len(stageHistoryRef.Archives)
	if revisionsStartIndex < 0 {
		revisionsStartIndex = 0
	}
	if startIndex < len(stageHistoryRef.Archives)-1 {
		stageHistoryRef.Archives = stageHistoryRef.Archives[:startIndex+1]
	}

	stageHistoryRef.Revisions = append(stageHistoryRef.Revisions[:revisionsStartIndex+1], revisionIDs...)
	found := false
	for _, revision := range stageHistoryRef.Revisions {
		if revision == activeRevisionID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Cannot set active revision to '%s': no corresponding revision", activeRevisionID)
	}
	stageHistoryRef.ActiveRevision = activeRevisionID

	// Generate keys
	var keysToDelete []interface{}
	var msetKeys []interface{}

	for _, revisionID := range revisionsToDelete {
		keysToDelete = append(keysToDelete, fmt.Sprintf("stageRevision:%s", revisionID))
	}

	for _, revision := range revisions {
		var revisionJSON []byte
		revisionJSON, err = json.Marshal(revision)
		if err != nil {
			return err
		}

		msetKeys = append(msetKeys,
			fmt.Sprintf("stageRevision:%s", revision.ID),
			revisionJSON)
	}

	historyJSON, err := json.Marshal(stageHistoryRef)
	if err != nil {
		return err
	}
	msetKeys = append(msetKeys, fmt.Sprintf("stageHistory:%s", stageID), historyJSON)

	// Start transaction
	err = conn.Send("MULTI")
	if err != nil {
		panic(err)
	}

	if len(keysToDelete) > 0 {
		err = conn.Send("DEL", keysToDelete...)
		if err != nil {
			panic(err)
		}
	}

	if len(msetKeys) > 0 {
		err = conn.Send("MSET", msetKeys...)
		if err != nil {
			panic(err)
		}
	}

	// End transaction
	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

// UpdateStageHistoryActiveRevision updates a stage history active revision
func UpdateStageHistoryActiveRevision(stageID, activeRevisionID string) error {
	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", fmt.Sprintf("stageHistory:%s", stageID)))
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("No stage history for id '%s'", stageID)
	}

	stageHistoryKey := fmt.Sprintf("stageHistory:%s", stageID)
	stageHistoryJSON, err := redis.Bytes(conn.Do("GET", stageHistoryKey))
	if err != nil {
		return err
	}

	var stageHistoryRef model.StageHistoryRef
	err = json.Unmarshal(stageHistoryJSON, &stageHistoryRef)
	if err != nil {
		return err
	}

	revisionIndex := -1
	revisions := append(stageHistoryRef.Archives, stageHistoryRef.Revisions...)

	for i, revisionID := range revisions {
		if revisionID == activeRevisionID {
			revisionIndex = i
			break
		}
	}

	if revisionIndex < 0 {
		return fmt.Errorf("No revision with ID '%s'", activeRevisionID)
	}

	stageHistoryRef.ActiveRevision = activeRevisionID
	stageHistoryJSON, err = json.Marshal(stageHistoryRef)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", stageHistoryKey, stageHistoryJSON)
	if err != nil {
		return err
	}

	return nil
}
