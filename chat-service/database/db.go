package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error

	// Get the database host from environment variable (defaults to "127.0.0.1" for local)
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1" // Default for local development
	}

	// Format the DSN (Data Source Name)
	dsn := fmt.Sprintf("user:password@tcp(%s:3306)/elderly?parseTime=true&loc=Local", dbHost)

	// Open the database connection
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Ping the database to check if the connection is established
	if err = DB.Ping(); err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	log.Println("Chatbot Database connected successfully")

}
