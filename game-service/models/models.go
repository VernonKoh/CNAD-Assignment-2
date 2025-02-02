package models

import "time"

type GameScore struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Score     int       `json:"score"`
	Timestamp time.Time `json:"timestamp"`
}
