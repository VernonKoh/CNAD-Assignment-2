package models

import "time"

// User represents a user in the system
type User struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	Name              string    `json:"name"`
	Role              string    `json:"role"`
	CreatedAt         time.Time `json:"created_at"`
	IsVerified        bool      `json:"is_verified"`        // Add this field
	VerificationToken string    `json:"verification_token"` // Add this field if needed for verification handling
}

// Struct to hold user profile data
type UserProfile struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	HighRisk    bool   `json:"high_risk"` // New field
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}
