package API

import (
	"database/sql"
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

	// Successful response
	response := APIResponse{Status: http.StatusOK, Message: "Success", User: user}
	sendResponse(w, response)
}

func createUser(username, email, password, biography, profilePic string) error {
	if len(profilePic) == 0 {
		profilePic = "/static/images/userAvatar/default-user.png"
	}
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
	_, err = stmt.Exec(username, email, password, biography, profilePic)
	if err != nil {
		return err
	}
	return nil
}

func GetUsersFollowers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendResponse(w, APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}

	// Extract user_id from URL parameter
	vars := mux.Vars(r)
	userID := vars["id"]
	if userID == "" {
		sendResponse(w, APIResponse{Status: http.StatusBadRequest, Message: "Missing user_id parameter"})
		return
	}

	// Adjust your SQL query to select count of followers
	var followerCount int
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM follow
        WHERE user_id = ?
    `, userID).Scan(&followerCount)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}

	// Prepare response with the follower count
	response := APIResponse{
		Status:       http.StatusOK,
		Message:      "Success",
		FollowerData: map[string]int{"follower_count": followerCount},
	}
	sendResponse(w, response)
}

func GetUsersFollowing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendResponse(w, APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}

	// Extract user_id from URL parameter
	vars := mux.Vars(r)
	userID := vars["id"]
	if userID == "" {
		sendResponse(w, APIResponse{Status: http.StatusBadRequest, Message: "Missing user_id parameter"})
		return
	}

	// Adjust your SQL query to select count of followers
	var followingCount int
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM followuser
        WHERE user_id = ?
    `, userID).Scan(&followingCount)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}

	// Prepare response with the follower count
	response := APIResponse{
		Status:       http.StatusOK,
		Message:      "Success",
		FollowerData: map[string]int{"following_count": followingCount},
	}
	sendResponse(w, response)
}

func GetUsersFollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendResponse(w, APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}

	// Extract user_id from URL parameter
	vars := mux.Vars(r)
	userID := vars["id"]
	if userID == "" {
		sendResponse(w, APIResponse{Status: http.StatusBadRequest, Message: "Missing user_id parameter"})
		return
	}

	// Adjust your SQL query to select count of followers
	var followingCount int
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM followuser
        WHERE follower_id = ?
    `, userID).Scan(&followingCount)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}

	// Prepare response with the follower count
	response := APIResponse{
		Status:       http.StatusOK,
		Message:      "Success",
		FollowerData: map[string]int{"follower_count": followingCount},
	}
	sendResponse(w, response)
}

func IsFollower(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendResponse(w, APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}

	// Extract user_id from URL parameter
	vars := mux.Vars(r)
	userID := vars["myId"]
	otherID := vars["otherId"]
	if userID == "" {
		sendResponse(w, APIResponse{Status: http.StatusBadRequest, Message: "Missing user.user_id parameter"})
		return
	}
	if otherID == "" {
		sendResponse(w, APIResponse{Status: http.StatusBadRequest, Message: "Missing user_id parameter"})
		return
	}

	var following bool
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM followUser
        WHERE follower_id = ? AND user_id = ?
    `, userID, otherID).Scan(&following)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}
	// Prepare response with the follower count
	response := APIResponse{
		Status:     http.StatusOK,
		Message:    "Success",
		IsFollower: following,
	}
	sendResponse(w, response)
}
