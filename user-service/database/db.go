package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Global database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	// Load environment variables from .env file
	err := godotenv.Load() // No need for a specific path since it's in user-service/
	if err != nil {
		log.Println("⚠️ Warning: .env file not found, using default values")
	}

	// Get credentials from .env or use default values
	mysqlUser := os.Getenv("MYSQL_USER")
	if mysqlUser == "" {
		mysqlUser = "root"
	}
	mysqlPassword := os.Getenv("MYSQL_PASS")
	if mysqlPassword == "" {
		mysqlPassword = "yourpassword"
	}
	mysqlHost := os.Getenv("MYSQL_HOST")
	if mysqlHost == "" {
		mysqlHost = "127.0.0.1"
	}
	mysqlDatabase := os.Getenv("MYSQL_DBNAME")
	if mysqlDatabase == "" {
		mysqlDatabase = "elderly"
	}

	// Construct DSN (Database Source Name)
	dsn := mysqlUser + ":" + mysqlPassword + "@tcp(" + mysqlHost + ":3306)/" + mysqlDatabase

	// Open database connection
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to the database: %v", err)
	}

	// Verify the connection
	if err = DB.Ping(); err != nil {
		log.Fatalf("❌ Database connection error: %v", err)
	}

	log.Println("✅ Database connected successfully")
}

// GetDB returns the database connection instance
func GetDB() *sql.DB {
	return DB
}
