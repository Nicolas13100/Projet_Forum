package API

import (
	"fmt"
	"net/http"
	"strconv"
)

// LikeTopicHandler handles liking a topic
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

	// Check if the topicID exists and corresponds to an existing topic
	if !topicExists(topicID) {
		response := APIResponse{Status: http.StatusNotFound, Message: "Topic not found"}
		sendResponse(w, response)
		return
	}

	// Check if the user has already liked the topic
	if alreadyLikedTopic(userID, topicID) {
		err := neutralLikeTopic(userID, topicID)
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "Couldn't remove like on topic"}
			sendResponse(w, response)
			return
		}
		response := APIResponse{Status: http.StatusOK, Message: "Like removed successfully"}
		sendResponse(w, response)
		return
	}

	// Add a record to the react_topic table indicating that the user liked the topic
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

// DislikeTopicHandler handles disliking a topic
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

	// Check if the topicID exists and corresponds to an existing topic
	if !topicExists(topicID) {
		response := APIResponse{Status: http.StatusNotFound, Message: "Topic not found"}
		sendResponse(w, response)
		return
	}

	// Check if the user has already disliked the topic
	if alreadyDislikedTopic(userID, topicID) {
		err := neutralDislikeTopic(userID, topicID)
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "Couldn't remove dislike from topic"}
			sendResponse(w, response)
			return
		}
		response := APIResponse{Status: http.StatusOK, Message: "Dislike removed successfully"}
		sendResponse(w, response)
		return
	}

	// Add a record to the react_topic table indicating that the user disliked the topic
	err = dislikeTopic(userID, topicID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to dislike topic"}
		sendResponse(w, response)
		return
	}

	response := APIResponse{Status: http.StatusOK, Message: "Topic disliked successfully"}
	sendResponse(w, response)
}

// Helper function to check if a topic exists in the database
func topicExists(topicID int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Topics_Table WHERE topic_id = ?)", topicID).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return exists
}

// Helper function to check if a user has already liked this topic in the database
func alreadyLikedTopic(userID, topicID int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM react_topic WHERE user_id = ? AND topic_id = ? AND status = '1')", userID, topicID).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return exists
}

// Helper function to check if a user has already disliked this topic in the database
func alreadyDislikedTopic(userID, topicID int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM react_topic WHERE user_id = ? AND topic_id = ? AND status = '2')", userID, topicID).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return exists
}

// Helper function to add a record to react_topic table indicating that the user liked the topic
func likeTopic(userID, topicID int) error {
	_, err := db.Exec("INSERT INTO react_topic (user_id, topic_id, status) VALUES (?, ?, '1')", userID, topicID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Helper function to add a record to react_topic table indicating that the user disliked the topic
func dislikeTopic(userID, topicID int) error {
	_, err := db.Exec("INSERT INTO react_topic (user_id, topic_id, status) VALUES (?, ?, '2')", userID, topicID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Helper function to remove like from a topic
func neutralLikeTopic(userID, topicID int) error {
	_, err := db.Exec("DELETE FROM react_topic WHERE user_id = ? AND topic_id = ? AND status = '1'", userID, topicID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Helper function to remove dislike from a topic
func neutralDislikeTopic(userID, topicID int) error {
	_, err := db.Exec("DELETE FROM react_topic WHERE user_id = ? AND topic_id = ? AND status = '2'", userID, topicID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
