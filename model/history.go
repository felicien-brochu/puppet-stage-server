package model

import "time"

// StageRevision revision of a stage
type StageRevision struct {
	ID    string    `json:"id"`
	Date  time.Time `json:"date"`
	Stage Stage     `json:"stage"`
}

// StageHistoryRef history of a stage. Stores references to existing revisions
type StageHistoryRef struct {
	StageID        string   `json:"stageID"`
	ActiveRevision string   `json:"activeRevision"`
	Revisions      []string `json:"revisions"`
	Archives       []string `json:"archives"`
}

// StageHistory history of a stage. Stores revisions
type StageHistory struct {
	StageID        string          `json:"stageID"`
	ActiveRevision string          `json:"activeRevision"`
	Revisions      []StageRevision `json:"revisions"`
}
