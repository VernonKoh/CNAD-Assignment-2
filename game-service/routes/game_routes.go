package routes

import (
	"CNAD_Assignment_2/game-service/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterGameRoutes defines all game-related routes
func RegisterGameRoutes(router *mux.Router) {
	// ✅ Ensure this line is present
	router.HandleFunc("/game/submit", handlers.SubmitScore).Methods("POST")
	router.HandleFunc("/game/scores/{id}", handlers.GetUserScores).Methods("GET") // Fetch scores
	// ✅ ADD Health Check Route (NEW)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Game Service is healthy"))
	}).Methods("GET")

}

//http://localhost:8083/game/scores/4
