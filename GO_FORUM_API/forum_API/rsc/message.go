package API

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func GetAllTopicMessagesAndAnswers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendResponse(w, APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	topicID := vars["topicID"]

	// Fetch all messages for the topic
	messages, err := fetchMessages(topicID)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	// Build a map of messageID to *Message
	messageMap := make(map[int]*Message)
	for _, msg := range messages {
		messageMap[msg.MessageID] = msg
	}

	// Fetch replies for each message
	for _, message := range messages {
		fetchReplies(message, messageMap)
	}

	sendResponse(w, APIResponse{Status: http.StatusOK, Message: "Success", TopicMessages: messages})
}

func fetchMessages(topicID string) ([]*Message, error) {
	rows, err := db.Query(`
        SELECT message_id, body, date_sent, topic_id, user_id
        FROM Messages_Table
        WHERE topic_id = ?`, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var message Message
		var dateSentStr string

		if err := rows.Scan(&message.MessageID, &message.Body, &dateSentStr, &message.TopicID, &message.UserID); err != nil {
			return nil, err
		}
		dateSent, err := time.Parse("2006-01-02 15:04:05", dateSentStr)
		if err != nil {
			return nil, err
		}
		message.DateSent = dateSent

		if err := fetchUserDetails(&message); err != nil {
			return nil, err
		}

		msgCopy := message                    // Create a copy of the message
		messages = append(messages, &msgCopy) // Append the pointer
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func fetchReplies(parentMessage *Message, messageMap map[int]*Message) {
	answerRows, err := db.Query(`
        SELECT message_id, body, date_sent, base_message_id, user_id
        FROM Messages_Table
        WHERE base_message_id = ?`, parentMessage.MessageID)
	if err != nil {
		log.Printf("Error querying replies: %v", err)
		return
	}
	defer answerRows.Close()

	for answerRows.Next() {
		var reply Message
		var dateSentStr string

		if err := answerRows.Scan(&reply.MessageID, &reply.Body, &dateSentStr, &reply.BaseMessageID, &reply.UserID); err != nil {
			log.Printf("Error scanning reply row: %v", err)
			return
		}
		dateSent, err := time.Parse("2006-01-02 15:04:05", dateSentStr)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			return
		}
		reply.DateSent = dateSent

		if err := fetchUserDetails(&reply); err != nil {
			log.Printf("Error fetching user details: %v", err)
			return
		}

		// Append the reply to the parent message
		parentMessage.Replies = append(parentMessage.Replies, &reply)
		// Recursively fetch replies for this reply
		fetchReplies(&reply, messageMap)
	}

	if err := answerRows.Err(); err != nil {
		log.Printf("Error iterating over reply rows: %v", err)
	}
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
