package API

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
)

func banUserHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	adminID, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	// Check if the user making the request has admin rights
	isAdmin, err := checkAdminRights(adminID)
	if err != nil {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to check admin rights"}
		sendResponse(w, response)
		return
	}
	if !isAdmin {
		response := APIResponse{Status: http.StatusUnauthorized, Message: "Unauthorized: You do not have admin rights"}
		sendResponse(w, response)
		return
	}

	// Parse form data to get the user ID and reason for deletion
	err = r.ParseForm()
	if err != nil {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to parse form data"}
		sendResponse(w, response)
		return
	}

	deletedUserID, _ := strconv.Atoi(r.FormValue("userID"))
	reason := r.FormValue("reason")

	// Create a log entry for the ban action
	err = createBanLog(reason)
	if err != nil {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to create ban log"}
		sendResponse(w, response)
		return
	}

	// Update the user's data to mark it as deleted
	err = deleteUser(deletedUserID)
	if err != nil {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to delete user"}
		sendResponse(w, response)
		return
	}

	response := APIResponse{Status: http.StatusOK, Message: "User Deleted correctly user"}
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
// deleteUser marks a user as deleted (soft deletion) in the database
func deleteUser(userID int) error {
	// Execute the UPDATE statement to set is_deleted to 1 for the specified user ID
	_, err := db.Exec("UPDATE users_table SET is_deleted = 1 WHERE user_id = ?", userID)
	if err != nil {
		return err
	}
	return nil
}

func modifyTopicHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	adminID, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	// Retrieve topicID from the request
	topicID := r.FormValue("topicID")

	// Query the database to retrieve the authorID associated with the topicID
	var authorID int
	err := db.QueryRow("SELECT user_id FROM Topics_Table WHERE topic_id = ?", topicID).Scan(&authorID)
	if err != nil {
		// Handle the error (e.g., topic not found)
		response := APIResponse{Status: http.StatusNotFound, Message: "Topic not found"}
		sendResponse(w, response)
		return
	}

	// Check if the adminID is the same as the authorID
	isAdmin, err := checkAdminRights(adminID)
	if adminID != authorID || !isAdmin {
		response := APIResponse{Status: http.StatusUnauthorized, Message: "Unauthorized: Only the topic author can modify it"}
		sendResponse(w, response)
		return
	}

	// If adminID matches authorID, proceed with modification
	// Parse form values for potential updates
	title := r.FormValue("title")
	body := r.FormValue("body")
	status := r.FormValue("status")
	isPrivate := r.FormValue("is_private")

	// Execute the update query
	_, err = db.Exec("UPDATE Topics_Table SET title = ?, body = ?, status = ?, is_private = ? WHERE topic_id = ?", title, body, status, isPrivate, authorID)
	if err != nil {
		// Handle the error
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to update topic"}
		sendResponse(w, response)
		return
	}

	// Respond with success message or redirect as needed
	response := APIResponse{Status: http.StatusOK, Message: "Topic updated successfully"}
	sendResponse(w, response)
}
