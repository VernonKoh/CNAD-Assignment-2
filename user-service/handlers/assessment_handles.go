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
	SELECT q.id, q.question_text, q.type, o.id, o.option_text, o.risk_value
	FROM Questions q
	LEFT JOIN Options o ON q.id = o.question_id
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
		var qID, oID, risk int
		var qText, qType, oText string

		err := rows.Scan(&qID, &qText, &qType, &oID, &oText, &risk)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, exists := questionsMap[qID]; !exists {
			questionsMap[qID] = &Question{
				ID:           qID,
				QuestionText: qText,
				Type:         qType,
				Options:      []Option{},
			}
		}

		if oID != 0 {
			questionsMap[qID].Options = append(questionsMap[qID].Options, Option{ID: oID, OptionText: oText, RiskValue: risk})
		}
	}

	var questions []Question
	for _, q := range questionsMap {
		questions = append(questions, *q)
	}
	log.Printf("Fetched assessment: ID=%s", assessmentID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

// Assessment struct
type Assessment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetAssessments returns all available quizzes
func GetAssessments(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name FROM assessments")

	if err != nil {
		http.Error(w, "Failed to fetch assessments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var assessments []Assessment
	for rows.Next() {
		var a Assessment
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			http.Error(w, "Error scanning assessments", http.StatusInternalServerError)
			return
		}
		assessments = append(assessments, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assessments)
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
	result, err := database.DB.Exec("INSERT INTO CompletedAssessments (assessment_id, user_id, total_risk_score, completed_at) VALUES (?, ?, ?, ?)",
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
		_, err := database.DB.Exec("INSERT INTO SelectedOptions (completed_id, option_id) VALUES (?, ?)", completedID, optionID)
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
