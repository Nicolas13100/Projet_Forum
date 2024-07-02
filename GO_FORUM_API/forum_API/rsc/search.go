package API

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
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

	response := APIResponse{
		Status:        http.StatusOK,
		Message:       "Search results",
		SearchResults: results,
	}

	sendResponse(w, response)
}

// Helper functions to search topics and messages
func searchTopics(query string) ([]Topic, error) {
	queryString := "SELECT topic_id, title, body, creation_date, status, is_private, user_id FROM Topics_Table WHERE title LIKE ? OR body LIKE ?"
	rows, err := db.Query(queryString, "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

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
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	var messages []Message
	for rows.Next() {
		var m Message
		var baseMessageID sql.NullInt64
		var dateSentString string
		if err := rows.Scan(&m.MessageID, &m.Body, &dateSentString, &m.TopicID, &baseMessageID, &m.UserID); err != nil {
			return nil, err
		}

		// Parse the dateSentString into time.Time
		m.DateSent, err = time.Parse("2006-01-02 15:04:05", dateSentString)
		if err != nil {
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

func GetForYouUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendResponse(w, APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}

	// Adjust your SQL query to select 3 random users with user_id, username, and profile_pic
	rows, err := db.Query(`
        SELECT user_id, username, profile_pic
        FROM users_table
        ORDER BY RAND()
        LIMIT 3
    `)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}
	defer rows.Close()

	// Prepare a slice to hold the results
	var users []User // Assuming User is a struct containing user_id, username, and profile_pic

	// Iterate through the rows and populate the users slice
	for rows.Next() {
		var user User
		// Scan the row into the User struct fields
		if err := rows.Scan(&user.UserID, &user.Username, &user.ProfilePic); err != nil {
			log.Printf("Error scanning row: %v", err)
			sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
			return
		}
		// Append each user to the slice
		users = append(users, user)
	}

	// Check for errors during rows iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}

	// If no users were found (though this case is unlikely with LIMIT 3)
	if len(users) == 0 {
		sendResponse(w, APIResponse{Status: http.StatusNotFound, Message: "No users found"})
		return
	}

	// Send the users slice as JSON response
	sendResponse(w, APIResponse{
		Status:    http.StatusOK,
		Message:   "Success",
		UsersData: users,
	})
}
