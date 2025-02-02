package main

import (
	"fmt"
	"log"
	"net/http"

	"CNAD_Assignment_2/game-service/database"
	"CNAD_Assignment_2/game-service/routes"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize the database
	database.InitDB()

	// Create a new router
	router := mux.NewRouter()

	// Register game routes
	routes.RegisterGameRoutes(router)

	// Add CORS support
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	// Start the server
	port := ":8083"
	fmt.Printf("Game Service is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
