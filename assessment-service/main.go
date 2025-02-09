package main

import (
	"CNAD_Assignment_2/assessment-service/routes" // Import your routes package
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

	// Serve the openpose.html file directly
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../Frontend/assessment.html")
	})

	// Wrap the router with CORS middleware
	handler := enableCORS(r)

	// Start the server
	fmt.Println("Server is running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", handler))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Game Service is healthy"))
	})

}
