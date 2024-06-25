package API

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func BanUserHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	adminID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"})
		return
	}

	// Check if the user making the request has admin rights
	isAdmin, err := checkAdminRights(adminID)
	if err != nil {
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Failed to check admin rights"})
		return
	}
	if !isAdmin {
		sendResponse(w, APIResponse{Status: http.StatusUnauthorized, Message: "Unauthorized: You do not have admin rights"})
		return
	}

	// Parse form data to get the user ID and reason for ban
	err = r.ParseForm()
	if err != nil {
		sendResponse(w, APIResponse{Status: http.StatusBadRequest, Message: "Failed to parse form data"})
		return
	}

	userID := r.FormValue("userID")
	reason := r.FormValue("reason")

	// Validate user ID
	deletedUserID, err := strconv.Atoi(userID)
	if err != nil {
		sendResponse(w, APIResponse{Status: http.StatusBadRequest, Message: "Invalid userID format"})
		return
	}

	// Create a log entry for the ban action
	messageLog := fmt.Sprintf("Admin %d banned User %d. Reason: %s", adminID, deletedUserID, reason)
	err = createBanLog(messageLog)
	if err != nil {
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Failed to create ban log"})
		return
	}

	// Update the user's data to mark it as banned
	err = banUser(deletedUserID)
	if err != nil {
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Failed to ban user"})
		return
	}

	response := APIResponse{Status: http.StatusOK, Message: "User banned successfully"}
	sendResponse(w, response)
}

// Function to check if the user making the request has admin rights
func checkAdminRights(userID int) (bool, error) {
	// Query the database to check if the user has admin rights
	var isAdmin bool
	err := db.QueryRow("SELECT isAdmin FROM users_table WHERE user_id = ?", userID).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // User not found, assuming no admin rights
		}
		return false, err
	}
	return isAdmin, nil
}

// createBanLog creates a log entry for the ban action
func createBanLog(reason string) error {
	// Insert a log entry into the Admin_Logs_Table
	_, err := db.Exec("INSERT INTO Admin_Logs_Table (action_descrbition) VALUES (?)", reason)
	if err != nil {
		return err
	}
	return nil
}

// Function to update the user's data to mark it as deleted
// banUser marks a user as deleted (soft deletion) in the database
func banUser(userID int) error {
	// Execute the UPDATE statement to set is_deleted to 1 for the specified user ID
	_, err := db.Exec("UPDATE users_table SET is_deleted = 1 WHERE user_id = ?", userID)
	if err != nil {
		return err
	}
	return nil
}

func ModifyTopicHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	adminID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"})
		return
	}

	// Retrieve topicID from the request
	topicID := r.FormValue("topicID")
	if topicID == "" {
		sendResponse(w, APIResponse{Status: http.StatusBadRequest, Message: "topicID is required"})
		return
	}

	// Query the database to retrieve the authorID associated with the topicID
	var authorID int
	err := db.QueryRow("SELECT user_id FROM Topics_Table WHERE topic_id = ?", topicID).Scan(&authorID)
	if err != nil {
		var response APIResponse
		if errors.Is(err, sql.ErrNoRows) {
			response = APIResponse{Status: http.StatusNotFound, Message: "Topic not found"}
		} else {
			log.Printf("Error querying topic: %v", err)
			response = APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
		}
		sendResponse(w, response)
		return
	}

	// Check if the user has admin rights
	isAdmin, err := checkAdminRights(adminID)
	if err != nil {
		log.Printf("Error checking admin rights: %v", err)
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}

	// Check if the adminID is the same as the authorID or the user is an admin
	if adminID != authorID && !isAdmin {
		sendResponse(w, APIResponse{Status: http.StatusUnauthorized, Message: "Unauthorized: Only the topic author or an admin can modify it"})
		return
	}

	// Parse form values for potential updates
	title := r.FormValue("title")
	body := r.FormValue("body")
	status := r.FormValue("status")
	isPrivate := r.FormValue("is_private")

	// Validate form values
	if title == "" || body == "" || status == "" || isPrivate == "" {
		sendResponse(w, APIResponse{Status: http.StatusBadRequest, Message: "All fields (title, body, status, is_private) are required"})
		return
	}

	// Execute the update query
	_, err = db.Exec("UPDATE Topics_Table SET title = ?, body = ?, status = ?, is_private = ? WHERE topic_id = ?", title, body, status, isPrivate, topicID)
	if err != nil {
		log.Printf("Error updating topic: %v", err)
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Failed to update topic"})
		return
	}

	// Respond with success message
	sendResponse(w, APIResponse{Status: http.StatusOK, Message: "Topic updated successfully"})
}
