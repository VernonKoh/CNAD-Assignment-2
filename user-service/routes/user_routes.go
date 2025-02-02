package routes

import (
	"CNAD_Assignment_2/user-service/handlers"

	"github.com/gorilla/mux"
)

// RegisterUserRoutes sets up the API routes for User Management
func RegisterUserRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/api/v1/users").Subrouter()
	userRouter.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	userRouter.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	userRouter.HandleFunc("/verify", handlers.VerifyUser).Methods("GET") // Add this route for verification

	// Add doctor login route
	userRouter.HandleFunc("/doctor/login", handlers.DoctorLogin).Methods("POST")

	// ✅ Add route for updating Facial ID
	userRouter.HandleFunc("/update-facial-id", handlers.UpdateFacialID).Methods("POST")

	// Add the webhook route to handle FACEIO events
	userRouter.HandleFunc("/faceio-webhook", handlers.HandleFaceIOWebhook).Methods("POST")

}

//get user profile
//http://localhost:8081/api/v1/users/2

//update user profile
//curl -X PUT "http://localhost:8081/api/v1/users/user_profile/2" -H "Content-Type: application/json" -d "{\"email\": \"updated.email@example.com\", \"name\": \"Updated Name\", \"age\": 65, \"gender\": \"Male\", \"address\": \"123 Updated Address, City\", \"phone_number\": \"1234567890\"}"
