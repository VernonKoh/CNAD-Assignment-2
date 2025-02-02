package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
<<<<<<< HEAD

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
=======
	DB, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/elderly?parseTime=true&loc=Local")
>>>>>>> a892430054abc094c0711c32ca3f644029111df5
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	log.Println("Database connected successfully")
}
