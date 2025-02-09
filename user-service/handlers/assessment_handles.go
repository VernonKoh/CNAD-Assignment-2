package handlers

import (
	"CNAD_Assignment_2/user-service/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Question represents a single question with options
type Question struct {
	ID           int      `json:"id"`
	QuestionText string   `json:"question_text"`
	Type         string   `json:"type"`
	Options      []Option `json:"options"`
}

// Option represents a single answer option
type Option struct {
	ID         int    `json:"id"`
	OptionText string `json:"option_text"`
	RiskValue  int    `json:"risk_value"`
}

// Fetch all questions and their options
func GetQuestions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Extract URL parameters
	assessmentID, ok := vars["assessment_id"]
	if !ok {
		http.Error(w, "Missing assessment ID", http.StatusBadRequest)
		return
	}

	// Prepare the SQL query with dynamic assessment_id
	query := `
	SELECT q.id, q.question_text, q.type, o.id AS option_id, o.option_text, o.risk_value
	FROM questions q
	LEFT JOIN options o ON q.id = o.question_id
	WHERE q.assessment_id = ?
	ORDER BY q.id ASC, o.id ASC;
	`

	rows, err := database.DB.Query(query, assessmentID)
	if err != nil {
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	questionsMap := make(map[int]*Question)

	for rows.Next() {
		var qID int
		var qText, qType string
		var oID *int      // Use a pointer to allow for NULL values in option ID
		var oText *string // Use a pointer to handle NULL in option_text
		var risk *int     // Use a pointer to handle NULL in risk_value

		// Scan the values, where oID is a pointer to int and oText is a pointer to string
		err := rows.Scan(&qID, &qText, &qType, &oID, &oText, &risk)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// If this is a new question, initialize it
		if _, exists := questionsMap[qID]; !exists {
			questionsMap[qID] = &Question{
				ID:           qID,
				QuestionText: qText,
				Type:         qType,
				Options:      []Option{}, // Initialize the Options slice
			}
		}

		// Only append options if oID and oText are not nil
		if oID != nil && oText != nil {
			questionsMap[qID].Options = append(questionsMap[qID].Options, Option{
				ID:         *oID,               // Dereference the pointer to get the ID
				OptionText: *oText,             // Dereference the pointer to get the option text
				RiskValue:  getRiskValue(risk), // Use helper function to get the value of risk
			})
		}
	}

	// Convert the map to a slice
	var questions []Question
	for _, q := range questionsMap {
		questions = append(questions, *q)
	}

	log.Printf("Fetched assessment: ID=%s", assessmentID)

	// Return the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

// Helper function to handle possible NULL values for risk_value
func getRiskValue(risk *int) int {
	if risk == nil {
		return 0 // Default value when risk is NULL
	}
	return *risk
}

// Assessment struct
type Assessment struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetAssessments returns all available quizzes
func GetAssessments(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, description FROM assessments")
	if err != nil {
		http.Error(w, "Failed to fetch assessments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var assessments []Assessment
	for rows.Next() {
		var a Assessment
		if err := rows.Scan(&a.ID, &a.Name, &a.Description); err != nil {
			http.Error(w, "Error scanning assessments", http.StatusInternalServerError)
			return
		}
		assessments = append(assessments, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assessments)
}

// GetAssessmentByID returns a single assessment based on the ID in the URL path
func GetAssessmentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get path parameters
	id := vars["assessmentid"]

	var a Assessment
	err := database.DB.QueryRow("SELECT id, name FROM assessments WHERE id = ?", id).Scan(&a.ID, &a.Name)

	if err != nil {
		http.Error(w, "Assessment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

type Submission struct {
	AssessmentID int   `json:"assessment_id"`
	UserID       int   `json:"user_id"`
	TotalRisk    int   `json:"total_risk"`
	OptionIDs    []int `json:"option_ids"`
}

// Save assessment submission
func SubmitAssessment(w http.ResponseWriter, r *http.Request) {
	var submission Submission
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Step 1: Insert into CompletedAssessments
	result, err := database.DB.Exec("INSERT INTO completed_assessments (assessment_id, user_id, total_risk_score, completed_at) VALUES (?, ?, ?, ?)",
		submission.AssessmentID, submission.UserID, submission.TotalRisk, time.Now())
	if err != nil {
		http.Error(w, "Failed to insert assessment", http.StatusInternalServerError)
		return
	}

	completedID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get last inserted ID", http.StatusInternalServerError)
		return
	}

	// Step 2: Insert into SelectedOptions
	for _, optionID := range submission.OptionIDs {
		_, err := database.DB.Exec("INSERT INTO selected_options (completed_id, option_id) VALUES (?, ?)", completedID, optionID)
		if err != nil {
			http.Error(w, "Failed to insert selected options", http.StatusInternalServerError)
			return
		}
	}

	// Success response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, `{"message": "Assessment submitted successfully"}`)
	log.Printf("Assessment submitted successfully")
}

func UpdateHighRiskStatus(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request struct {
		UserID   int  `json:"user_id"`
		HighRisk bool `json:"high_risk"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// Update high risk status
	query := "UPDATE users SET high_risk = ? WHERE id = ?"
	_, err := database.DB.Exec(query, request.HighRisk, request.UserID)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User high risk status updated successfully"})
	log.Printf("User ID: %d high risk status updated successfully to %t", request.UserID, request.HighRisk)

}
