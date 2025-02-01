package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// Global database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	var err error

	// Hardcoded database connection string (not recommended for production)
	DB, err = sql.Open("mysql", "root:yourpassword@tcp(127.0.0.1:3306)/elderly")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Verify the connection
	if err = DB.Ping(); err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	log.Println("✅ Database connected successfully")
}

// GetDB returns the database connection instance
func GetDB() *sql.DB {
	return DB
}
