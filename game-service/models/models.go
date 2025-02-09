package models

import "time"

// ScoreRequest struct
type ScoreRequest struct {
	UserID    int `json:"user_id"`
	Score     int `json:"score"`
	TimeTaken int `json:"time_taken"`
}

// ScoreResponse struct for returning scores
type ScoreResponse struct {
	Score     int       `json:"score"`
	TimeTaken int       `json:"time_taken"`
	Timestamp time.Time `json:"timestamp"` // New field to hold the timestamp
}
