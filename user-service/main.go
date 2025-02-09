package main

import (
	"CNAD_Assignment_2/user-service/database"
	"CNAD_Assignment_2/user-service/handlers"
	"CNAD_Assignment_2/user-service/notification"
	"CNAD_Assignment_2/user-service/routes"

	"fmt"
	"log"
	"net/http"
	"time"

	gorillaHandlers "github.com/gorilla/handlers" // ✅ Rename to avoid conflict
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database
	database.InitDB()

	// Run the check immediately on startup
	log.Println("Running initial high-risk check...")
	notification.NotifyUsers()

	// Schedule the function to run every hour
	ticker := time.NewTicker(1 * time.Hour) // Runs every 1 hour
	defer ticker.Stop()

	go func() { // ✅ Run in a separate goroutine
		for {
			<-ticker.C // Wait for the next tick (1 hour)
			log.Println("Running scheduled high-risk check...")
			notification.NotifyUsers()
		}
	}()

	// Create a new router
	r := mux.NewRouter()

	// Register API routes for user management
	routes.RegisterUserRoutes(r)

	// ✅ Test route to check if API is working
	r.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is running"))
	}).Methods("GET")

	// Update Facial ID
	r.HandleFunc("/api/v1/users/update-facial-id", handlers.UpdateFacialID).Methods("POST")

	// Serve the static files for the frontend
	staticDir := "../frontend" // Path to the directory containing `index.html`
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	// ✅ Enable CORS
	corsMiddleware := gorillaHandlers.CORS( // ✅ Use `gorillaHandlers` instead of `handlers`
		gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Start the server with CORS enabled
	fmt.Println("Server is running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", corsMiddleware(r)))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User Service is healthy"))
	})
}
