package handlers

import (
	"CNAD_Assignment_2/game-service/database"
	"encoding/json"
	"log"
	"net/http"

	"time"

	"github.com/gorilla/mux"
)

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

// ✅ Function to handle score submission
func SubmitScore(w http.ResponseWriter, r *http.Request) {
	var request ScoreRequest

	// Decode JSON request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error": "Invalid request data"}`, http.StatusBadRequest)
		return
	}

	// Validate user_id, score, and time_taken
	if request.UserID == 0 || request.Score < 0 || request.TimeTaken < 0 {
		http.Error(w, `{"error": "Invalid input values"}`, http.StatusBadRequest)
		return
	}

	// ✅ Insert game results into the database
	query := `INSERT INTO game_scores (user_id, score, time_taken) VALUES (?, ?, ?)`
	_, err = database.DB.Exec(query, request.UserID, request.Score, request.TimeTaken)
	if err != nil {
		log.Println("❌ Error inserting score:", err)
		http.Error(w, `{"error": "Failed to save score"}`, http.StatusInternalServerError)
		return
	}

	// ✅ Success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Score submitted successfully!"})
}

// ✅ Function to retrieve user game scores
func GetUserScores(w http.ResponseWriter, r *http.Request) {
	// Get user ID from request URL
	vars := mux.Vars(r)
	userID := vars["id"]

	// Query the database for the user's game scores, including timestamp
	rows, err := database.DB.Query("SELECT score, time_taken, timestamp FROM game_scores WHERE user_id = ? ORDER BY time_taken ASC", userID)
	if err != nil {
		log.Println("❌ Error fetching scores:", err)
		http.Error(w, `{"error": "Failed to retrieve scores"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Store retrieved scores
	var scores []ScoreResponse
	for rows.Next() {
		var score ScoreResponse
		if err := rows.Scan(&score.Score, &score.TimeTaken, &score.Timestamp); err != nil {
			log.Println("❌ Error scanning score row:", err)
			http.Error(w, `{"error": "Failed to parse scores"}`, http.StatusInternalServerError)
			return
		}
		scores = append(scores, score)
	}

	// ✅ Handle case where no scores exist for user
	if len(scores) == 0 {
		log.Println("⚠️ No scores found for user", userID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]ScoreResponse{}) // ✅ Return an empty array instead of null
		return
	}

	// ✅ Send JSON response with scores
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scores)
}
