package main

import (
	"CNAD_Assignment_2/assessment-service/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Register API routes for assessment (OpenPose API)
	routes.RegisterAssessmentRoutes(r)

	// Serve static files from the ../Frontend/ directory
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../Frontend/")))

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"},                   // Allow requests from this origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allow GET and POST methods
		AllowedHeaders:   []string{"Authorization", "Content-Type"},           // Allow these headers
		AllowCredentials: true,                                                // Allow credentials (e.g., cookies)
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(r)

	// Start the server
	fmt.Println("Server is running on http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", handler))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Game Service is healthy"))
	})
}
