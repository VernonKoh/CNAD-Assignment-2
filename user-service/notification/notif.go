package main

import (
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv" // Import the package for loading .env files
)

func main() {
	// Load environment variables from .env file (if present)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	// Get SMTP credentials from environment variables
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Check if required environment variables are set
	if smtpUser == "" || smtpPass == "" {
		log.Fatal("SMTP_USER or SMTP_PASS environment variable not set")
	}

	// Email details
	from := smtpUser
	to := "curtislee1028@gmail.com" // Replace with the recipient's email
	subject := "Test Email"
	body := "Hello! This is a secure test email sent from Go."

	// Compose the email message
	message := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	// Set up authentication information
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Send the email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Email sent successfully!")
}
