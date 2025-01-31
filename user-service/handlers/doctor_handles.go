package handlers

import (
	"CNAD_Assignment_2/user-service/database"
	"CNAD_Assignment_2/user-service/models"
	"CNAD_Assignment_2/user-service/utils"
	"encoding/json"
	"log"
	"net/http"
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
