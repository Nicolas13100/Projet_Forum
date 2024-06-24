package API

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/forumcour")
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
