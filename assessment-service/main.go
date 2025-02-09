package main

import (
	"CNAD_Assignment_2/assessment-service/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// enableCORS adds CORS headers to the response
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow specific methods
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Register API routes for assessment (OpenPose API)
	routes.RegisterAssessmentRoutes(r)

	// Serve static files from the ../Frontend/ directory
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../Frontend/")))

	// Wrap the router with CORS middleware
	handler := enableCORS(r)

	// Start the server
	fmt.Println("Server is running on http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", handler))
}
