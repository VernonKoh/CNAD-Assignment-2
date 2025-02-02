package handlers

import (
	"CNAD_Assignment_2/user-service/database"
	"CNAD_Assignment_2/user-service/models"
	"CNAD_Assignment_2/user-service/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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

// FacialIDUpdateRequest struct
type FacialIDUpdateRequest struct {
	UserID   int    `json:"userID"`
	FacialID string `json:"facialID"`
}

// UpdateFacialID updates the user's facial ID in the database
func UpdateFacialID(w http.ResponseWriter, r *http.Request) {
	// ✅ Ensure request is JSON
	w.Header().Set("Content-Type", "application/json")

	var request FacialIDUpdateRequest

	// ✅ Read request body once
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "Failed to read request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// ✅ Decode JSON properly
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// ✅ Validate input
	if request.UserID == 0 || request.FacialID == "" {
		http.Error(w, `{"error": "Missing userID or facialID"}`, http.StatusBadRequest)
		return
	}

	// ✅ Update Facial ID in database (Fixed Query)
	query := "UPDATE users SET facial_id = ? WHERE id = ?"
	_, err = database.DB.Exec(query, request.FacialID, request.UserID)
	if err != nil {
		http.Error(w, `{"error": "Database update failed"}`, http.StatusInternalServerError)
		return
	}

	// ✅ Success Response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Facial ID updated successfully"})
}

// FaceIOWebhookPayload represents the expected structure of the FACEIO webhook request
type FaceIOWebhookPayload struct {
	Event    string `json:"event"`
	UserID   int    `json:"userID"`
	FacialID string `json:"facialID"`
}

// HandleFaceIOWebhook processes webhook events from FACEIO
func HandleFaceIOWebhook(w http.ResponseWriter, r *http.Request) {
	var payload FaceIOWebhookPayload

	// ✅ Set response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// ✅ Read and log the raw request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("❌ Error reading request body:", err)
		http.Error(w, `{"error": "Failed to read request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Println("📩 Raw Webhook Body:", string(body))

	// ✅ Decode the JSON request body
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Println("❌ Invalid JSON payload:", err)
		http.Error(w, `{"error": "Invalid JSON data"}`, http.StatusBadRequest)
		return
	}

	log.Printf("✅ Parsed Webhook Data: %+v\n", payload)

	// ✅ Process different FACEIO events
	switch payload.Event {
	case "faceio.enrollment.success":
		log.Println("🔹 Face enrollment successful for user:", payload.UserID)

		// ✅ Update facial ID in the database
		query := "UPDATE users SET facial_id = ? WHERE id = ?"
		result, err := database.DB.Exec(query, payload.FacialID, payload.UserID)
		if err != nil {
			log.Printf("❌ SQL Error updating facial ID: %v", err)
			http.Error(w, `{"error": "Database update failed"}`, http.StatusInternalServerError)
			return
		}

		// ✅ Check if any rows were affected
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			log.Printf("⚠️ No rows updated! UserID %d might not exist.", payload.UserID)
			http.Error(w, `{"error": "User not found or already updated"}`, http.StatusBadRequest)
			return
		}

		log.Println("✅ Facial ID updated successfully for user:", payload.UserID)

	case "faceio.authentication.success":
		log.Println("🔹 User authenticated successfully with facial recognition:", payload.FacialID)

		// ✅ Verify if facial ID exists in the database
		var userID int
		err := database.DB.QueryRow("SELECT id FROM users WHERE facial_id = ?", payload.FacialID).Scan(&userID)
		if err != nil {
			log.Printf("❌ Facial ID not found in database: %v", err)
			http.Error(w, `{"error": "Facial ID not recognized"}`, http.StatusUnauthorized)
			return
		}

		log.Printf("✅ User with Facial ID %s authenticated: UserID %d", payload.FacialID, userID)

	case "faceio.facialid.delete":
		log.Println("❌ Facial ID deletion requested for user:", payload.UserID)

		// ✅ Remove Facial ID from Database
		query := "UPDATE users SET facial_id = NULL WHERE id = ?"
		result, err := database.DB.Exec(query, payload.UserID)
		if err != nil {
			log.Printf("❌ Error deleting facial ID: %v", err)
			http.Error(w, `{"error": "Database update failed"}`, http.StatusInternalServerError)
			return
		}

		// ✅ Check if any rows were affected
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			log.Printf("⚠️ No rows deleted! UserID %d might not exist.", payload.UserID)
			http.Error(w, `{"error": "User not found or already deleted"}`, http.StatusBadRequest)
			return
		}

		log.Println("🗑️ Facial ID removed successfully for user:", payload.UserID)

	default:
		log.Println("⚠️ Unhandled FACEIO event:", payload.Event)
	}

	// ✅ Send Success Response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Webhook processed successfully"})
}

// GetEmailByFaceID retrieves the email associated with a given facial ID
func GetEmailByFaceID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request struct {
		FacialID string `json:"facialID"`
	}

	// ✅ Read and parse request body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("❌ Invalid request payload:", err)
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// ✅ Validate input
	if request.FacialID == "" {
		log.Println("❌ Missing facial ID")
		http.Error(w, `{"error": "Facial ID is required"}`, http.StatusBadRequest)
		return
	}

	// ✅ Query the database to get the user ID and email
	var userID int
	var userEmail string

	query := "SELECT id, email FROM users WHERE facial_id = ?"
	err = database.DB.QueryRow(query, request.FacialID).Scan(&userID, &userEmail)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("⚠️ No user found for facialID: %s", request.FacialID)
		} else {
			log.Printf("❌ Database error: %v", err)
		}
		http.Error(w, `{"error": "Facial ID not found"}`, http.StatusNotFound)
		return
	}

	// ✅ Send back the email and user ID in JSON format
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"userID": userID,
		"email":  userEmail,
	})

}

// UserDetailsResponse struct
type UserDetailsResponse struct {
	UserID int    `json:"userID"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}

// GetUserDetails returns user details based on email
func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	var userDetails UserDetailsResponse
	err = database.DB.QueryRow("SELECT id, name, email FROM users WHERE email = ?", request.Email).
		Scan(&userDetails.UserID, &userDetails.Name, &userDetails.Email)

	if err != nil {
		log.Printf("❌ Error fetching user details: %v", err)
		http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
		return
	}

	// ✅ Generate a JWT token (for session persistence)
	userDetails.Token, err = utils.GenerateJWT(userDetails.UserID)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userDetails)
}
