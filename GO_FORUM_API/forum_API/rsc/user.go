package API

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		response := APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"}
		sendResponse(w, response)
		return
	}

	// Extract user_id from URL parameter
	vars := mux.Vars(r)
	userID := vars["id"]

	// Query the database for the user with the given ID
	user := User{}
	err := db.QueryRow(`
        SELECT user_id, username, email, 
               registration_date, last_login_date, biography, 
               isAdmin, isModerator, is_deleted, profile_pic 
        FROM users_table 
        WHERE user_id = ?`, userID).
		Scan(&user.UserID, &user.Username, &user.Email,
			&user.RegistrationDate, &user.LastLoginDate, &user.Biography,
			&user.IsAdmin, &user.IsModerator, &user.IsDeleted, &user.ProfilePic)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no user found with the given ID, return 404 Not Found
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			// For other errors, return 500 Internal Server Error
			log.Println("Error querying database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Convert user struct to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshaling user to JSON:", err)
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	response := APIResponse{Status: http.StatusOK, Message: "Success", JsonResp: userJSON}
	sendResponse(w, response)
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
