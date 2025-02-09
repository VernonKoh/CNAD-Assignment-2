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

	// Serve the static files for the frontend
	staticDir := "../frontend" // Path to the directory containing `index.html`
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:8081", "http://localhost:8083"}, // Allow both 8081 and 8083
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	// Start the server
	port := ":8083"
	fmt.Printf("Game Service is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Game Service is healthy"))
	})

}
