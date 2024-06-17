package API

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func GetAllTopicMessage(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		response := APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"}
		sendResponse(w, response)
		return
	}

	// Extract user_id from URL parameter
	vars := mux.Vars(r)
	topicID := vars["id"]

	// Query the database for the user with the given ID
	user := User{}
	err := db.QueryRow(`
        SELECT *
        FROM messages_table 
        WHERE topic_id = ?`, topicID).
		Scan(&user.UserID, &user.Username, &user.Email,
			&user.RegistrationDate, &user.LastLoginDate, &user.Biography,
			&user.IsAdmin, &user.IsModerator, &user.IsDeleted, &user.ProfilePic)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no user found with the given ID, return 404 Not Found
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			// For other errors, return 500 Internal Server Error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Convert user struct to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	response := APIResponse{Status: http.StatusOK, Message: "Success", JsonResp: userJSON}
	sendResponse(w, response)
}
