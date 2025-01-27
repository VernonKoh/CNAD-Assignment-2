package main

import (
	"CNAD_Assignment_2/user-service/database"
	"CNAD_Assignment_2/user-service/notification"
	"CNAD_Assignment_2/user-service/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database
	database.InitDB()

	// Call the notification function
	notification.NotifyUsers()

	log.Println("Application has finished processing notifications.")

	// Create a new router
	r := mux.NewRouter()

	// Register API routes for user management
	routes.RegisterUserRoutes(r)

	// Serve the static files for the frontend
	staticDir := "../frontend" // Path to the directory containing `index.html`
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	// Start the server
	fmt.Println("Server is running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
