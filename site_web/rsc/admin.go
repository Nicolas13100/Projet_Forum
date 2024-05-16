package API

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func banUserHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	adminID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Check if the user making the request has admin rights
	isAdmin, err := checkAdminRights(adminID)
	if err != nil {
		http.Error(w, "Failed to check admin rights", http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "Unauthorized: You do not have admin rights", http.StatusUnauthorized)
		return
	}

	// Parse form data to get the user ID and reason for deletion
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
		return
	}

	deletedUserID, _ := strconv.Atoi(r.FormValue("userID"))
	reason := r.FormValue("reason")

	// Create a log entry for the ban action
	err = createBanLog(reason)
	if err != nil {
		http.Error(w, "Failed to create ban log", http.StatusInternalServerError)
		return
	}

	// Update the user's data to mark it as deleted
	err = deleteUser(deletedUserID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Optionally, you can send a success response
	w.WriteHeader(http.StatusOK)
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
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Retrieve topicID from the request
	topicID := r.FormValue("topicID")

	// Query the database to retrieve the authorID associated with the topicID
	var authorID int
	err := db.QueryRow("SELECT user_id FROM Topics_Table WHERE topic_id = ?", topicID).Scan(&authorID)
	if err != nil {
		// Handle the error (e.g., topic not found)
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}

	// Check if the adminID is the same as the authorID
	isAdmin, err := checkAdminRights(adminID)
	if adminID != authorID || !isAdmin {
		http.Error(w, "Unauthorized: Only the topic author can modify it", http.StatusUnauthorized)
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
		http.Error(w, "Failed to update topic", http.StatusInternalServerError)
		return
	}

	// Respond with success message or redirect as needed
	_, err = w.Write([]byte("Topic updated successfully"))
	if err != nil {
		fmt.Println(err)
		return
	}
}
