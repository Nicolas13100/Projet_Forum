package API

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func likeCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	// Retrieve messageID from the request
	messageIDStr := r.FormValue("messageID")
	if messageIDStr == "" {
		response := APIResponse{Status: http.StatusBadRequest, Message: "messageID parameter is required"}
		sendResponse(w, response)
		return
	}

	// Convert messageID from string to int
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		response := APIResponse{Status: http.StatusBadRequest, Message: "Invalid messageID"}
		sendResponse(w, response)
		return
	}

	// Check if the messageID exists and corresponds to an existing message
	if !messageExists(messageID) {
		http.Error(w, "Message not found", http.StatusNotFound)
		response := APIResponse{Status: http.StatusNotFound, Message: "Message not found"}
		sendResponse(w, response)
		return
	}

	// Check if the user has already liked the comment
	if alreadyLiked(userID, messageID) {
		err := neutralLikeComment(userID, messageID)
		if err != nil {
			response := APIResponse{Status: http.StatusInternalServerError, Message: "Couldn't take off like"}
			sendResponse(w, response)
			return
		}
		return
	}

	// Add a record to the react_message table indicating that the user liked the comment
	err = likeComment(userID, messageID)
	if err != nil {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to like comment"}
		sendResponse(w, response)
		return
	}

	// Create API response
	response := Response{
		Message: "Comment liked successfully",
	}

	// Encode response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to encode response"}
		sendResponse(w, response)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to write response"}
		sendResponse(w, response)
		return
	}
}

func dislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	// Retrieve messageID from the request
	messageIDStr := r.FormValue("messageID")
	if messageIDStr == "" {
		response := APIResponse{Status: http.StatusBadRequest, Message: "messageID parameter is required"}
		sendResponse(w, response)
		return
	}

	// Convert messageID from string to int
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		response := APIResponse{Status: http.StatusBadRequest, Message: "Invalid messageID"}
		sendResponse(w, response)
		return
	}

	// Check if the messageID exists and corresponds to an existing message
	if !messageExists(messageID) {
		response := APIResponse{Status: http.StatusNotFound, Message: "Message not found"}
		sendResponse(w, response)
		return
	}

	// Check if the user has already liked the comment
	if alreadyDisliked(userID, messageID) {
		err := neutralLikeComment(userID, messageID)
		if err != nil {
			response := APIResponse{Status: http.StatusInternalServerError, Message: "couldn't take off dislike"}
			sendResponse(w, response)
			return
		}
		return
	}

	// Add a record to the react_message table indicating that the user liked the comment
	err = dislikeComment(userID, messageID)
	if err != nil {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to like comment"}
		sendResponse(w, response)
		return
	}

	// Create API response
	response := Response{
		Message: "Comment liked successfully",
	}

	// Encode response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to encode response"}
		sendResponse(w, response)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to write response"}
		sendResponse(w, response)
		return
	}
}

// Helper function to check if a message exists in the database
func messageExists(messageID int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Messages_Table WHERE message_id = ?)", messageID).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return exists
}

// Helper function to check if a user has already liked a comment in the database
func alreadyLiked(userID, messageID int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM react_message WHERE user_id = ? AND message_id = ?)", userID, messageID).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return exists
}

// Helper function to check if a user has already liked a comment in the database
func alreadyDisliked(userID, messageID int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 2 FROM react_message WHERE user_id = ? AND message_id = ?)", userID, messageID).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return exists
}

// Helper function to add a record to react_message table indicating that the user liked the comment
func likeComment(userID, messageID int) error {
	_, err := db.Exec("INSERT INTO react_message (user_id, message_id, status) VALUES (?, ?, 1)", userID, messageID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Helper function to add a record to react_message table indicating that the user liked the comment
func dislikeComment(userID, messageID int) error {
	_, err := db.Exec("INSERT INTO react_message (user_id, message_id, status) VALUES (?, ?, 2)", userID, messageID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Helper function to add a record to react_message table indicating that the user liked the comment
func neutralLikeComment(userID, messageID int) error {
	_, err := db.Exec("INSERT INTO react_message (user_id, message_id, status) VALUES (?, ?, 0)", userID, messageID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
