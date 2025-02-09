package routes

import (
	"CNAD_Assignment_2/assessment-service/handlers" // Import your handlers package

	"github.com/gorilla/mux"
)

// RegisterAssessmentRoutes sets up the API routes for OpenPose processing
func RegisterAssessmentRoutes(router *mux.Router) {
	assessmentRouter := router.PathPrefix("/api/v1/assessment").Subrouter()
	assessmentRouter.HandleFunc("/upload", handlers.UploadHandler).Methods("POST") // Route for file upload
	assessmentRouter.HandleFunc("/upload_video", handlers.UploadVideoHandler).Methods("POST")
}
