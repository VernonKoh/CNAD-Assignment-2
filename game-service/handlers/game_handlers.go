package handlers

import (
	"CNAD_Assignment_2/game-service/database"
	"encoding/json"
	"log"
	"net/http"
)

// ScoreRequest struct to handle JSON data from frontend
type ScoreRequest struct {
	UserID    int `json:"user_id"`
	Score     int `json:"score"`
	TimeTaken int `json:"time_taken"`
}

// SubmitScore saves the game result in the database
func SubmitScore(w http.ResponseWriter, r *http.Request) {
	var request ScoreRequest

	// Decode JSON request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error": "Invalid request data"}`, http.StatusBadRequest)
		return
	}

	// Insert score into database
	query := `INSERT INTO game_scores (user_id, score, time_taken) VALUES (?, ?, ?)`
	_, err = database.DB.Exec(query, request.UserID, request.Score, request.TimeTaken)
	if err != nil {
		log.Println("Error inserting score:", err)
		http.Error(w, `{"error": "Failed to save score"}`, http.StatusInternalServerError)
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Score submitted successfully"})
}
