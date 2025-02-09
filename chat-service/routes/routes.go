package routes

import (
	"CNAD_Assignment_2/chat-service/handlers" // ✅ Corrected import

	"github.com/gorilla/mux"
)

// SetupRoutes registers all API routes
func SetupRoutes(router *mux.Router) {
	router.HandleFunc("/chat", handlers.ChatHandler).Methods("POST")
}
