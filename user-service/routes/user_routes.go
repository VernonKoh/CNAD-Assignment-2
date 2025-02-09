package routes

import (
	"CNAD_Assignment_2/user-service/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterUserRoutes sets up the API routes for User Management
func RegisterUserRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/api/v1/users").Subrouter()
	userRouter.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	userRouter.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	userRouter.HandleFunc("/verify", handlers.VerifyUser).Methods("GET") // Add this route for verification

	userRouter.HandleFunc("/user_profile/{id}", handlers.GetUserProfile).Methods("GET")
	userRouter.HandleFunc("/user_profile/{id}", handlers.UpdateUserProfile).Methods("PUT")

	// Add doctor login route
	userRouter.HandleFunc("/doctor/login", handlers.DoctorLogin).Methods("POST")

	// ✅ Add route for updating Facial ID
	userRouter.HandleFunc("/update-facial-id", handlers.UpdateFacialID).Methods("POST")

	userRouter.HandleFunc("/get-email-by-faceid", handlers.GetEmailByFaceID).Methods("POST")

	userRouter.HandleFunc("/get-user-details", handlers.GetUserDetails).Methods("POST")

	// Add the webhook route to handle FACEIO events
	userRouter.HandleFunc("/faceio-webhook", handlers.HandleFaceIOWebhook).Methods("POST")

	userRouter.HandleFunc("/questions/{assessment_id}", handlers.GetQuestions).Methods("GET")
	userRouter.HandleFunc("/assessments", handlers.GetAssessments).Methods("GET")
	userRouter.HandleFunc("/assessments/{assessmentid}", handlers.GetAssessmentByID).Methods("GET")

	userRouter.HandleFunc("/submit-assessment", handlers.SubmitAssessment).Methods("POST")
	userRouter.HandleFunc("/update-high-risk", handlers.UpdateHighRiskStatus).Methods("POST")

	//doctor functions
	userRouter.HandleFunc("/search", handlers.SearchUsers).Methods("GET")
	userRouter.HandleFunc("/completed_assessments/{userid}", handlers.GetCompletedAssessments).Methods("GET")
	userRouter.HandleFunc("/assessments", handlers.CreateAssessment).Methods("POST") // New route
	userRouter.HandleFunc("/questions", handlers.CreateQuestion).Methods("POST")
	userRouter.HandleFunc("/options", handlers.CreateOption).Methods("POST")
	userRouter.HandleFunc("/assessments/{id}", handlers.DeleteAssessment).Methods("DELETE")

	// ✅ ADD Health Check Route (NEW)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User Service is healthy"))
	}).Methods("GET")

}

//access mysql
//mysql -u root -p

//get user profile
//http://localhost:8081/api/v1/users/user_profile/2

//update user profile
//curl -X PUT "http://localhost:8081/api/v1/users/user_profile/2" -H "Content-Type: application/json" -d "{\"email\": \"updated.email@example.com\", \"name\": \"Updated Name\", \"age\": 65, \"gender\": \"Male\", \"address\": \"123 Updated Address, City\", \"phone_number\": \"1234567890\"}"

//get all questions and options for assessment
//curl -X GET http://localhost:8081/api/v1/users/questions/1
//curl -X GET http://localhost:8081/api/v1/users/assessments

//change user to high risk
//curl -X POST "http://localhost:8081/api/v1/users/update-high-risk" -H "Content-Type: application/json" -d "{\"user_id\": 4, \"high_risk\": true}"

//search for patients
//http://localhost:8081/api/v1/users/search?name=e

//Get all of User's completed assessments
//http://localhost:8081/api/v1/users/completed_assessments/4

//create new assessments
//curl -X POST http://localhost:8081/api/v1/users/assessments -H "Content-Type: application/json" -d "{\"name\":\"Risk Assessment\",\"description\":\"An assessment to evaluate risk levels.\"}"

//create question
//curl -X POST http://localhost:8081/api/v1/users/questions -H "Content-Type: application/json" -d "{\"assessment_id\":4,\"question_text\":\"What is your risk level?\",\"type\":\"mcq\"}"

//create option
//curl -X POST http://localhost:8081/api/v1/users/options -H "Content-Type: application/json" -d "{\"assessment_id\":4,\"question_id\":1,\"option_text\":\"Low Risk\",\"risk_value\":1}"

//delete assessment
//curl -X DELETE http://localhost:8081/api/v1/users/assessments/19
