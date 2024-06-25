package API

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var db *sql.DB

func InitDB() error {
	// Load the environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Build the data source name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		// Close the connection if the ping fails
		if closeErr := db.Close(); closeErr != nil {
			return fmt.Errorf("ping error: %v, close error: %v", err, closeErr)
		}
		return err
	}

	fmt.Println("Connection to MySQL database successfully established")
	return nil
}

func createUser(username, email, password, biography, profilePic string) error {
	hashedPassword := hashPassword(password)
	// Insert user into database
	stmt, err := db.Prepare("INSERT INTO users_table (username, email, password, biography, profile_pic) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(stmt)
	_, err = stmt.Exec(username, email, hashedPassword, biography, profilePic)
	if err != nil {
		return err
	}
	return nil
}
