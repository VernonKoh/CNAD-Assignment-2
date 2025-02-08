package handlers

import (
	"CNAD_Assignment_2/user-service/database"
	"CNAD_Assignment_2/user-service/models"
	"CNAD_Assignment_2/user-service/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// DoctorLogin handles login for doctors
func DoctorLogin(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
		return
	}

	// Fetch doctor from the database
	var doctor models.User
	query := "SELECT id, name, password, is_verified FROM doctors WHERE email = ?"
	err := database.DB.QueryRow(query, credentials.Email).Scan(&doctor.ID, &doctor.Name, &doctor.Password, &doctor.IsVerified)
	if err != nil {
		log.Printf("Error fetching doctor: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}

	// Check if doctor is verified
	if !doctor.IsVerified {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email not verified"})
		return
	}

	// Validate password (use direct comparison if plain text)
	if credentials.Password != doctor.Password {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(doctor.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate token"})
		return
	}

	// Respond with token, userID, and name
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":  token,
		"userID": doctor.ID,
		"name":   doctor.Name,
	})
}

// SearchUsers handles the search request for users by name
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameter 'name'
	query := r.URL.Query().Get("name")
	if query == "" {
		http.Error(w, "Query parameter 'name' is required", http.StatusBadRequest)
		return
	}

	// Query the database for users with names matching the search term
	rows, err := database.DB.Query(`
		SELECT id, name, high_risk
		FROM users
		WHERE name LIKE ?
	`, "%"+query+"%")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		HighRisk bool   `json:"high_risk"`
	}

	// Loop through the rows and append to the users slice
	for rows.Next() {
		var user struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			HighRisk bool   `json:"high_risk"`
		}
		err := rows.Scan(&user.ID, &user.Name, &user.HighRisk)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Return the list of users as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Struct for Completed Assessment
type CompletedAssessment struct {
	ID             int    `json:"id"`
	AssessmentID   int    `json:"assessment_id"`
	UserID         int    `json:"user_id"`
	TotalRiskScore int    `json:"total_risk_score"`
	CompletedAt    string `json:"completed_at"`
}

// Handler to fetch completed assessments by user_id
func GetCompletedAssessments(w http.ResponseWriter, r *http.Request) {
	// Get user_id from URL query params
	userID := mux.Vars(r)["userid"]
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Query the database
	query := "SELECT id, assessment_id, user_id, total_risk_score, completed_at FROM completedassessments WHERE user_id = ?"
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Store results in a slice
	var assessments []CompletedAssessment
	for rows.Next() {
		var assessment CompletedAssessment
		if err := rows.Scan(&assessment.ID, &assessment.AssessmentID, &assessment.UserID, &assessment.TotalRiskScore, &assessment.CompletedAt); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		assessments = append(assessments, assessment)
	}

	// Convert to JSON and send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assessments)
}

// Assessment struct
type NewAssessment struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateAssessment handles creating a new assessment
func CreateAssessment(w http.ResponseWriter, r *http.Request) {
	var newAssessment NewAssessment

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&newAssessment); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Insert into the database
	query := "INSERT INTO assessments (name, description) VALUES (?, ?)"
	result, err := database.DB.Exec(query, newAssessment.Name, newAssessment.Description)
	if err != nil {
		http.Error(w, `{"error": "Failed to create assessment"}`, http.StatusInternalServerError)
		return
	}

	// Get the newly inserted assessment ID
	assessmentID, _ := result.LastInsertId()

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Assessment created successfully",
		"assessment_id": assessmentID,
	})
}

// CreateQuestion handles creating a new question
func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON body
	var question struct {
		AssessmentID int    `json:"assessment_id"`
		QuestionText string `json:"question_text"`
		Type         string `json:"type"`
	}

	// Decode the request body into the question structure
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&question); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert the new question into the database
	query := `INSERT INTO questions (assessment_id, question_text, type) VALUES (?, ?, ?)`
	_, err := database.DB.Exec(query, question.AssessmentID, question.QuestionText, question.Type)
	if err != nil {
		http.Error(w, "Failed to create question", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Question created successfully"})
}

type newOption struct {
	ID           int    `json:"id"`
	AssessmentID int    `json:"assessment_id"`
	QuestionID   int    `json:"question_id"`
	OptionText   string `json:"option_text"`
	RiskValue    int    `json:"risk_value"`
}

// CreateOption creates a new option in the database
func CreateOption(w http.ResponseWriter, r *http.Request) {
	var newOption newOption

	// Parse the JSON body into the Option struct
	if err := json.NewDecoder(r.Body).Decode(&newOption); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// SQL query to insert the new option into the options table
	query := "INSERT INTO options (assessment_id, question_id, option_text, risk_value) VALUES (?, ?, ?, ?)"

	// Execute the query
	result, err := database.DB.Exec(query, newOption.AssessmentID, newOption.QuestionID, newOption.OptionText, newOption.RiskValue)
	if err != nil {
		http.Error(w, "Failed to create option", http.StatusInternalServerError)
		return
	}

	// Get the ID of the newly created option
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve inserted option ID", http.StatusInternalServerError)
		return
	}

	// Respond with the newly created option
	newOption.ID = int(lastInsertID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newOption)
}
