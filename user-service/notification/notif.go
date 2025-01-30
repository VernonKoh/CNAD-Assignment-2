package notification

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"CNAD_Assignment_2/user-service/database"

	"github.com/joho/godotenv" // Import for loading .env files
)

func NotifyUsers() {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	// Get SMTP credentials from environment variables
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	if smtpUser == "" || smtpPass == "" {
		log.Fatal("SMTP_USER or SMTP_PASS environment variable not set")
	}

	// Email setup
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Ensure database is initialized
	if database.DB == nil {
		log.Fatal("Database connection is not initialized. Call InitDB first.")
	}

	// High-Risk Alerts
	queryHighRisk := `
		SELECT u.email, a.risk_level, a.assessment_date
		FROM assessments a
		JOIN users u ON a.user_id = u.id
		WHERE a.risk_level = 'High'
	`
	rows, err := database.DB.Query(queryHighRisk)
	if err != nil {
		log.Fatalf("Failed to query database for high-risk users: %v", err)
	}
	defer rows.Close()

	// Send high-risk alerts
	for rows.Next() {
		var email string
		var riskLevel string
		var assessmentDate string

		if err := rows.Scan(&email, &riskLevel, &assessmentDate); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		subject := "Health Alert: High Risk Level Detected"
		body := fmt.Sprintf("Dear Senior,\n\nYour recent health assessment on %s flagged your risk level as 'High'. We recommend consulting a clinic for further evaluation.\n\nBest regards,\nHealthCare Team", assessmentDate)

		message := []byte("From: " + smtpUser + "\r\n" +
			"To: " + email + "\r\n" +
			"Subject: " + subject + "\r\n\r\n" +
			body + "\r\n")

		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{email}, message)
		if err != nil {
			log.Printf("Failed to send email to %s: %v", email, err)
		} else {
			log.Printf("Alert sent to %s for high risk level on %s", email, assessmentDate)

			// Testing Log: Print email sent details
			fmt.Printf("TEST LOG: Email sent successfully to %s with subject: %s\n", email, subject)
		}
	}

	log.Println("All high-risk alerts processed.")

	// Check-Up Reminders (Commented Out)
	/*
		now := time.Now()
		sevenDaysLater := now.AddDate(0, 0, 7).Format("2006-01-02")
		thirtyDaysLater := now.AddDate(0, 0, 30).Format("2006-01-02")

		queryCheckups := `
			SELECT email, next_checkup_date
			FROM checkup_schedules
			WHERE next_checkup_date BETWEEN ? AND ?
		`
		rowsCheckups, err := database.DB.Query(queryCheckups, sevenDaysLater, thirtyDaysLater)
		if err != nil {
			log.Fatalf("Failed to query database for check-up reminders: %v", err)
		}
		defer rowsCheckups.Close()

		// Send check-up reminders
		for rowsCheckups.Next() {
			var email string
			var checkupDate string

			if err := rowsCheckups.Scan(&email, &checkupDate); err != nil {
				log.Printf("Failed to scan row: %v", err)
				continue
			}

			subject := "Reminder: Upcoming Health Check-Up"
			body := fmt.Sprintf("Dear Senior,\n\nThis is a reminder that your next health check-up is scheduled for %s. Please make the necessary arrangements.\n\nBest regards,\nHealthCare Team", checkupDate)

			message := []byte("From: " + smtpUser + "\r\n" +
				"To: " + email + "\r\n" +
				"Subject: " + subject + "\r\n\r\n" +
				body + "\r\n")

			err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{email}, message)
			if err != nil {
				log.Printf("Failed to send email to %s: %v", email, err)
			} else {
				log.Printf("Reminder sent to %s for check-up on %s", email, checkupDate)

				// Testing Log: Print email sent details
				fmt.Printf("TEST LOG: Email sent successfully to %s with subject: %s\n", email, subject)
			}
		}

		log.Println("All check-up reminders processed.")
	*/
}
