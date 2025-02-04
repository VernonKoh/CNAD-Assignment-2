package handlers

import (
	"CNAD_Assignment_2/user-service/database"
	"encoding/json"
	"log"
	"net/http"

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
