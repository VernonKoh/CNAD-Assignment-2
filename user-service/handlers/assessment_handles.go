package handlers

import (
	"CNAD_Assignment_2/user-service/database"
	"encoding/json"
	"net/http"
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
	rows, err := database.DB.Query(`
		SELECT q.id, q.question_text, q.type, o.id, o.option_text, o.risk_value
		FROM Questions q
		LEFT JOIN Options o ON q.id = o.question_id
		WHERE q.assessment_id = 1
		ORDER BY q.id ASC, o.id ASC;
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}
