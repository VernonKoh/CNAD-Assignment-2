package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// Global database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	var err error

	// Read MySQL credentials from environment variables, fallback to default values
	mysqlUser := os.Getenv("MYSQL_USER")
	if mysqlUser == "" {
		mysqlUser = "root"
	}

	mysqlPassword := os.Getenv("MYSQL_PASS")
	if mysqlPassword == "" {
		mysqlPassword = "password" // Default password for others
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

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to the database: %v", err)
	}

	// Verify connection
	if err = DB.Ping(); err != nil {
		log.Fatalf("❌ Database connection error: %v", err)
	}

	log.Println("✅ Database connected successfully")
}

// GetDB returns the database connection instance
func GetDB() *sql.DB {
	return DB
}

// run this in cmd before starting go app (set up environment variables)
// set MYSQL_USER=root
// set MYSQL_PASS=yourpassword  # Replace with your own MySQL password
// set MYSQL_HOST=127.0.0.1
// set MYSQL_DBNAME=elderly
// go run main.go
