package routes

import (
	"CNAD_Assignment_2/game-service/handlers"

	"github.com/gorilla/mux"
)

// RegisterGameRoutes defines all game-related routes
func RegisterGameRoutes(router *mux.Router) {
	router.HandleFunc("/game/submit", handlers.SubmitScore).Methods("POST") // ✅ New route

}
