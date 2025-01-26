package notification

import (
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv" // Import for loading .env files
)

func main() {
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

	// Example email details
	email := "curtislee1028@gmail.com" // Replace with a valid email
	subject := "Test Email"
	body := "This is a test email to verify the email functionality."

	// Compose the email
	message := []byte("From: " + smtpUser + "\r\n" +
		"To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	// Send the email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{email}, message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
	} else {
		log.Printf("Email sent successfully to %s", email)
	}

	// Commented out database-related sections
	/*
		// Database connection (replace with your MySQL DSN)
		dsn := "your_user:your_password@tcp(localhost:3306)/your_db"
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		// Query seniors with upcoming check-ups (7-30 days from today)
		now := time.Now()
		sevenDaysLater := now.AddDate(0, 0, 7).Format("2006-01-02")
		thirtyDaysLater := now.AddDate(0, 0, 30).Format("2006-01-02")

		query := `
			SELECT email, next_checkup_date
			FROM checkup_schedules
			WHERE next_checkup_date BETWEEN ? AND ?
		`
		rows, err := db.Query(query, sevenDaysLater, thirtyDaysLater)
		if err != nil {
			log.Fatalf("Failed to query database: %v", err)
		}
		defer rows.Close()

		// Email sending loop (commented out for now)
		for rows.Next() {
			var email string
			var checkupDate string

			if err := rows.Scan(&email, &checkupDate); err != nil {
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
			}
		}

		log.Println("All reminders processed.")
	*/
}
