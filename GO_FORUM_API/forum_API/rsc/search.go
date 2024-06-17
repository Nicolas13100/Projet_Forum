package API

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// SearchHandler handles search queries for topics and messages
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	topics, err := searchTopics(query)
	if err != nil {
		log.Printf("Error searching topics: %v", err)
		http.Error(w, "Error searching topics", http.StatusInternalServerError)
		return
	}

	messages, err := searchMessages(query)
	if err != nil {
		log.Printf("Error searching messages: %v", err)
		http.Error(w, "Error searching messages", http.StatusInternalServerError)
		return
	}

	results := SearchResults{
		Topics:   topics,
		Messages: messages,
	}

	jsonResp, err := json.Marshal(results)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	response := APIResponse{
		Status:   http.StatusOK,
		Message:  "Search results",
		JsonResp: json.RawMessage(jsonResp),
	}

	sendResponse(w, response)
}

func searchTopics(query string) ([]Topic, error) {
	queryString := "SELECT topic_id, title, body, creation_date, status, is_private, user_id FROM Topics_Table WHERE title LIKE ? OR body LIKE ?"
	rows, err := db.Query(queryString, "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var t Topic
		if err := rows.Scan(&t.TopicID, &t.Title, &t.Body, &t.CreationDate, &t.Status, &t.IsPrivate, &t.UserID); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return topics, nil
}

func searchMessages(query string) ([]Message, error) {
	queryString := "SELECT message_id, body, date_sent, topic_id, base_message_id, user_id FROM Messages_Table WHERE body LIKE ?"
	rows, err := db.Query(queryString, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		var baseMessageID sql.NullInt64
		if err := rows.Scan(&m.MessageID, &m.Body, &m.DateSent, &m.TopicID, &baseMessageID, &m.UserID); err != nil {
			return nil, err
		}
		if baseMessageID.Valid {
			baseMessageIDValue := int(baseMessageID.Int64)
			m.BaseMessageID = &baseMessageIDValue
		}
		messages = append(messages, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
