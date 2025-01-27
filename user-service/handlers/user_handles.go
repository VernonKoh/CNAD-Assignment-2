package handlers

import (
	"CNAD_Assignment_2/user-service/database"
	"CNAD_Assignment_2/user-service/models"
	"CNAD_Assignment_2/user-service/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	log.Println("RegisterUser handler called")

	var user models.User

	// Decode JSON input
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding input: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
		return
	}

	// Debugging received data
	log.Printf("Received user data: Email=%s, Name=%s, Password=%s", user.Email, user.Name, user.Password)

	// Validate email
	if !utils.ValidateEmail(user.Email) {
		log.Println("Invalid email format")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email address"})
		return
	}

	// Validate password
	if strings.TrimSpace(user.Password) == "" {
		log.Println("Password is required")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Password is required"})
		return
	}

	// Validate name
	if strings.TrimSpace(user.Name) == "" {
		log.Println("Name is required")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Name is required"})
		return
	}

	// Hash the password
	log.Println("Hashing the password")
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error hashing password"})
		return
	}
	log.Printf("Hashed password: %s", hashedPassword)

	// Assign hashed password to the user
	user.Password = hashedPassword

	// Generate a verification token
	log.Println("Generating verification token")
	token, err := utils.GenerateVerificationToken()
	if err != nil {
		log.Printf("Error generating verification token: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error generating verification token"})
		return
	}

	// Start a transaction to ensure both inserts succeed or fail together
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	// Insert the user into the users table
	query := "INSERT INTO users (email, password, name, role, verification_token) VALUES (?, ?, ?, ?, ?)"
	result, err := tx.Exec(query, user.Email, user.Password, user.Name, user.Role, token)
	if err != nil {
		tx.Rollback()
		log.Printf("Error inserting user into database: %v", err)
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(err.Error(), "Duplicate entry") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "Email already registered"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to register user"})
		}
		return
	}
	// Get the inserted user's ID
	userID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Printf("Error getting last inserted user ID: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	// Insert into user_details with default values
	query = "INSERT INTO user_details (user_id, age, gender, address, phone_number) VALUES (?, ?, ?, ?, ?)"
	_, err = tx.Exec(query, userID, 0, "Unknown", "", "")
	if err != nil {
		tx.Rollback()
		log.Printf("Error inserting user details: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	// Simulate sending the email
	verificationLink := fmt.Sprintf("http://localhost:8081/api/v1/users/verify?token=%s", token)
	log.Printf("Send email to %s with verification link: %s", user.Email, verificationLink)

	// Send the actual verification email
	if err := utils.SendVerificationEmail(user.Email, verificationLink); err != nil {
		log.Printf("Error sending verification email: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to send verification email"})
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully. Please verify your email.",
		"email":   user.Email,
		"token":   token, // Optional: Include token for testing purposes
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
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

	// Fetch user from the database
	var user models.User
	query := "SELECT id, name, password, is_verified FROM users WHERE email = ?"
	err := database.DB.QueryRow(query, credentials.Email).Scan(&user.ID, &user.Name, &user.Password, &user.IsVerified)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}
	log.Printf("Fetched user: ID=%d, Name=%s, Password=%s, IsVerified=%t", user.ID, user.Name, user.Password, user.IsVerified)
	// Check if user is verified
	if !user.IsVerified {
		log.Printf("Email not verified for user: %s", credentials.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email not verified"})
		return
	}
	log.Printf("Login successful: UserID=%d, Name=%s", user.ID, user.Name)
	// Validate the password
	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		log.Printf("Invalid password for user: %s", credentials.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
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
		"userID": user.ID,
		"name":   user.Name,
	})
}

func VerifyUser(w http.ResponseWriter, r *http.Request) {
	// Extract the token from the query parameters
	token := r.URL.Query().Get("token")
	if token == "" {
		log.Println("Missing token in query string")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing token"})
		return
	}

	// Debug: Log the received token
	log.Printf("Verification token received: %s", token)

	// Check if the token exists in the database
	var userID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE verification_token = ?", token).Scan(&userID)
	if err != nil {
		log.Printf("Error finding user with token: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token or user not found"})
		return
	}

	// Debug: Log the user ID associated with the token
	log.Printf("User ID associated with token: %d", userID)

	// Update the user's verification status
	result, err := database.DB.Exec("UPDATE users SET is_verified = TRUE, verification_token = NULL WHERE id = ?", userID)
	if err != nil {
		log.Printf("Error updating user verification status: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error verifying user"})
		return
	}

	// Check if any row was affected
	rowsAffected, _ := result.RowsAffected()
	log.Printf("Rows affected: %d", rowsAffected)
	if rowsAffected == 0 {
		log.Println("Invalid token or user not found in database")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token or user not found"})
		return
	}

	log.Println("User email verified successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Email verified successfully"})
}

func UpdateUserMembership(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"] // Get the user ID from the URL
	if id == "" {
		log.Println("Error: Missing or invalid user ID")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}

	var request struct {
		Role string `json:"role"` // The new membership tier
	}

	log.Printf("Received membership update request: ID=%s", id) // Debug: Log ID

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request body: %v", err) // Debug: Log error
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
		return
	}

	// Validate membership role
	validRoles := map[string]bool{"Basic": true, "Premium": true, "VIP": true}
	if !validRoles[request.Role] {
		log.Printf("Invalid membership role: %s", request.Role)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid membership tier"})
		return
	}

	// Convert ID to integer
	userID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Error converting user ID to integer: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}

	log.Printf("Updating membership for user ID=%d to Role=%s", userID, request.Role) // Debug: Log role

	// Update the membership tier in the database
	result, err := database.DB.Exec("UPDATE users SET role = ? WHERE id = ?", request.Role, userID)
	if err != nil {
		log.Printf("Database error during update: %v", err) // Debug: Log error
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update membership tier"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("No rows updated for user ID: %d", userID)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found or no change in role"})
		return
	}

	log.Printf("Membership tier updated successfully for user ID=%d", userID) // Debug: Confirm success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Membership tier updated successfully", "updatedRole": request.Role})
}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"] // Extract user ID from the URL

	// Validate if ID is numeric
	userID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Invalid user ID: %s", id)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}

	var profile models.UserProfile

	// Fetch user data and details from the database
	query := `
		SELECT 
			u.id, u.email, u.name, u.role, 
			ud.age, ud.gender, ud.address, ud.phone_number
		FROM users u
		LEFT JOIN user_details ud ON u.id = ud.user_id
		WHERE u.id = ?`
	err = database.DB.QueryRow(query, userID).Scan(
		&profile.ID, &profile.Email, &profile.Name, &profile.Role,
		&profile.Age, &profile.Gender, &profile.Address, &profile.PhoneNumber,
	)
	if err != nil {
		log.Printf("Error fetching user profile for ID=%d: %v", userID, err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}

	log.Printf("Fetched user profile for ID=%d", userID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

// UpdateUserProfile allows users to update their details and membership
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"] // Extract user ID from URL
	// Validate if ID is numeric
	userID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Invalid user ID: %s", id)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}
	// Decode the JSON payload
	var updatedProfile models.UserProfile
	if err := json.NewDecoder(r.Body).Decode(&updatedProfile); err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	// Validate input fields (e.g., ensure no empty fields)
	if updatedProfile.Name == "" || updatedProfile.Email == "" {
		log.Println("Name or Email cannot be empty")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Name and Email are required"})
		return
	}
	// Update the `users` table
	userQuery := `
		UPDATE users 
		SET email = ?, name = ? 
		WHERE id = ?
	`
	_, err = database.DB.Exec(userQuery, updatedProfile.Email, updatedProfile.Name, userID)
	if err != nil {
		log.Printf("Error updating users table for ID=%d: %v", userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update user information"})
		return
	}

	// Update the `user_details` table
	detailsQuery := `
		UPDATE user_details 
		SET age = ?, gender = ?, address = ?, phone_number = ? 
		WHERE user_id = ?
	`
	_, err = database.DB.Exec(detailsQuery, updatedProfile.Age, updatedProfile.Gender, updatedProfile.Address, updatedProfile.PhoneNumber, userID)
	if err != nil {
		log.Printf("Error updating user_details table for user_id=%d: %v", userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update user details"})
		return
	}

	log.Printf("Updated user ID=%s: Name=%s, Age=%d", id, updatedProfile.Name, updatedProfile.Age)
	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User profile updated successfully"})
}
