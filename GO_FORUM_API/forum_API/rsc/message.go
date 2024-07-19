package API

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetAllTopicMessagesAndAnswers(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		sendResponse(w, APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}

	// Extract topic_id from URL parameter
	vars := mux.Vars(r)
	topicID := vars["topicID"]

	// Query the database for messages with the given topic_id
	rows, err := db.Query(`
        SELECT message_id, body, date_sent, topic_id, user_id
        FROM Messages_Table
        WHERE topic_id = ?`, topicID)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}
	defer rows.Close()

	// Iterate over the rows and build the list of messages
	var messages []Message
	for rows.Next() {
		var message Message
		var dateSentBytes []byte

		if err := rows.Scan(&message.MessageID, &message.Body, &dateSentBytes, &message.TopicID, &message.UserID); err != nil {
			log.Printf("Error scanning row: %v", err)
			sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
			return
		}

		dateSentStr := string(dateSentBytes)
		dateSent, err := time.Parse("2006-01-02 15:04:05", dateSentStr)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
			return
		}
		message.DateSent = dateSent

		// Fetch user details
		if err := fetchUserDetails(&message); err != nil {
			log.Printf("Error fetching user details: %v", err)
			sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
			return
		}

		messages = append(messages, message)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}

	// Recursively fetch answers for each message
	var allMessages []Message
	for _, message := range messages {
		allMessages = append(allMessages, message)
		fetchAnswers(&allMessages, strconv.Itoa(message.MessageID))
	}

	// Send the messages as JSON response
	sendResponse(w, APIResponse{Status: http.StatusOK, Message: "Success", TopicMessages: allMessages})
}

func fetchUserDetails(message *Message) error {
	row := db.QueryRow(`
        SELECT username, profile_pic
        FROM users_table
        WHERE user_id = ?`, message.UserID)

	var username, profilePic string
	if err := row.Scan(&username, &profilePic); err != nil {
		return err
	}

	message.Username = username
	message.ProfilePic = profilePic
	return nil
}

func fetchAnswers(allMessages *[]Message, messageID string) {
	answerRows, err := db.Query(`
        SELECT message_id, body, date_sent, base_message_id, user_id
        FROM Messages_Table
        WHERE base_message_id = ?`, messageID)
	if err != nil {
		log.Printf("Error querying answers: %v", err)
		return
	}
	defer answerRows.Close()

	// Iterate over the answer rows and add them to the list
	for answerRows.Next() {
		var answer Message
		var dateSentBytes []byte

		if err := answerRows.Scan(&answer.MessageID, &answer.Body, &dateSentBytes, &answer.BaseMessageID, &answer.UserID); err != nil {
			log.Printf("Error scanning answer row: %v", err)
			return
		}

		dateSentStr := string(dateSentBytes)
		dateSent, err := time.Parse("2006-01-02 15:04:05", dateSentStr)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			return
		}
		answer.DateSent = dateSent

		// Fetch user details
		if err := fetchUserDetails(&answer); err != nil {
			log.Printf("Error fetching user details: %v", err)
			return
		}

		*allMessages = append(*allMessages, answer)
		// Recursively fetch answers for this answer
		fetchAnswers(allMessages, strconv.Itoa(answer.MessageID))
	}

	// Check for errors from iterating over answer rows
	if err := answerRows.Err(); err != nil {
		log.Printf("Error iterating over answer rows: %v", err)
	}
}
