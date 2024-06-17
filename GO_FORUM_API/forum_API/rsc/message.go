package API

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func GetAllTopicMessage(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		sendResponse(w, APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}

	// Extract topic_id from URL parameter
	vars := mux.Vars(r)
	topicID := vars["id"]

	// Query the database for messages with the given topic_id
	rows, err := db.Query(`
        SELECT message_id, body, date_sent, topic_id, user_id
        FROM Messages_Table
        WHERE topic_id = ?`, topicID)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	// Iterate over the rows and build the list of messages
	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.MessageID, &message.Body, &message.DateSent, &message.TopicID, &message.UserID); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert messages slice to JSON
	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Send the messages as JSON response
	sendResponse(w, APIResponse{Status: http.StatusOK, Message: "Success", JsonResp: messagesJSON})
}

func GetAllMessageAnswer(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		sendResponse(w, APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}

	// Extract topic_id from URL parameter
	vars := mux.Vars(r)
	messageID := vars["id"]

	// Query the database for messages with the given topic_id
	rows, err := db.Query(`
        SELECT message_id, body, date_sent, base_message_id, user_id
        FROM Messages_Table
        WHERE message_id = ?`, messageID)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	// Iterate over the rows and build the list of messages
	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.MessageID, &message.Body, &message.DateSent, &message.BaseMessageID, &message.UserID); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert messages slice to JSON
	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Send the messages as JSON response
	sendResponse(w, APIResponse{Status: http.StatusOK, Message: "Success", JsonResp: messagesJSON})
}
