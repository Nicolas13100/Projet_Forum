package API

import (
	"fmt"
	"net/http"
	"strconv"
)

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println(err)
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
	if alreadyLiked(userID, messageID) {
		err := neutralLikeComment(userID, messageID)
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "Couldn't take off like"}
			sendResponse(w, response)
			return
		}
		return
	}

	// Add a record to the react_message table indicating that the user liked the comment
	err = likeComment(userID, messageID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to like comment"}
		sendResponse(w, response)
		return
	}
	response := APIResponse{Status: http.StatusOK, Message: "Comment liked successfully"}
	sendResponse(w, response)
}

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println(err)
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
			fmt.Println(err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "couldn't take off dislike"}
			sendResponse(w, response)
			return
		}
		return
	}

	// Add a record to the react_message table indicating that the user liked the comment
	err = dislikeComment(userID, messageID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to like comment"}
		sendResponse(w, response)
		return
	}

	response := APIResponse{Status: http.StatusOK, Message: "Comment disliked successfully"}
	sendResponse(w, response)
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

func LikeTopicHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	// Retrieve topicID from the request
	topicIDStr := r.FormValue("topicID")
	if topicIDStr == "" {
		response := APIResponse{Status: http.StatusBadRequest, Message: "topicID parameter is required"}
		sendResponse(w, response)
		return
	}

	// Convert topicID from string to int
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Invalid topicID"}
		sendResponse(w, response)
		return
	}

	// Check if the topicID exists and corresponds to an existing message
	if !topicExists(topicID) {
		response := APIResponse{Status: http.StatusNotFound, Message: "Message not found"}
		sendResponse(w, response)
		return
	}

	// Check if the user has already liked the comment
	if alreadyLikedTopic(userID, topicID) {
		err := neutralLikeTopic(userID, topicID)
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "Couldn't take off like on topic"}
			sendResponse(w, response)
			return
		}
		return
	}

	// Add a record to the react_topic table indicating that the user liked the comment
	err = likeTopic(userID, topicID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to like topic"}
		sendResponse(w, response)
		return
	}
	response := APIResponse{Status: http.StatusOK, Message: "Topic liked successfully"}
	sendResponse(w, response)
}

func DislikeTopicHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	// Retrieve topicID from the request
	topicIDStr := r.FormValue("topicID")
	if topicIDStr == "" {
		response := APIResponse{Status: http.StatusBadRequest, Message: "topicID parameter is required"}
		sendResponse(w, response)
		return
	}

	// Convert topicID from string to int
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Invalid topicID"}
		sendResponse(w, response)
		return
	}

	// Check if the topicID exists and corresponds to an existing message
	if !topicExists(topicID) {
		response := APIResponse{Status: http.StatusNotFound, Message: "topic not found"}
		sendResponse(w, response)
		return
	}

	// Check if the user has already liked the comment
	if alreadyDislikedTopic(userID, topicID) {
		err := neutralLikeComment(userID, topicID)
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "couldn't take off dislike from topic"}
			sendResponse(w, response)
			return
		}
		return
	}

	// Add a record to the react_topic table indicating that the user liked the comment
	err = dislikeTopic(userID, topicID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to like topic"}
		sendResponse(w, response)
		return
	}
	// Create API response
	response := APIResponse{Status: http.StatusOK, Message: "topic disliked successfully"}
	sendResponse(w, response)
}

// Helper function to check if a topic exists in the database
func topicExists(topicID int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM topics_table WHERE topic_id = ?)", topicID).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return exists
}

// Helper function to check if a user has already liked this topic in the database
func alreadyLikedTopic(userID, topicID int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM react_topic WHERE user_id = ? AND topic_id = ?)", userID, topicID).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return exists
}

// Helper function to check if a user has already disliked this topic in the database
func alreadyDislikedTopic(userID, topicID int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 2 FROM react_topic WHERE user_id = ? AND topic_id = ?)", userID, topicID).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return exists
}

// Helper function to add a record to react_topic table indicating that the user liked the comment
func likeTopic(userID, topicID int) error {
	_, err := db.Exec("INSERT INTO react_topic (user_id, topic_id, status) VALUES (?, ?, 1)", userID, topicID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Helper function to add a record to react_topic table indicating that the user disliked the comment
func dislikeTopic(userID, topicID int) error {
	_, err := db.Exec("INSERT INTO react_topic (user_id, topic_id, status) VALUES (?, ?, 2)", userID, topicID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Helper function to add a record to react_topic table indicating that the user already liked the comment
func neutralLikeTopic(userID, topicID int) error {
	_, err := db.Exec("INSERT INTO react_topic (user_id, topic_id, status) VALUES (?, ?, 0)", userID, topicID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
