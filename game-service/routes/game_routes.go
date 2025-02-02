package routes

import (
	"CNAD_Assignment_2/game-service/handlers"

	"github.com/gorilla/mux"
)

// RegisterGameRoutes defines all game-related routes
func RegisterGameRoutes(router *mux.Router) {
	// ✅ Ensure this line is present
	router.HandleFunc("/game/submit", handlers.SubmitScore).Methods("POST")
	router.HandleFunc("/game/scores/{id}", handlers.GetUserScores).Methods("GET") // Fetch scores
}
