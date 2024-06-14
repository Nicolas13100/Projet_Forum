package API

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func createTopicHandler(w http.ResponseWriter, r *http.Request) {

	//check method
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//parse data
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	//check if all the data are here
	title := r.FormValue("title")
	body := r.FormValue("body")
	status := r.FormValue("status")
	//check image
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	if title == "" || body == "" || status == "" || userID != 0 {
		http.Error(w, "Missing required form fields", http.StatusBadRequest)
		return
	}

	// Insert user into database
	stmt, err := db.Prepare("INSERT INTO Topics_Table (title, body, status, userID) VALUES (?, ?, ?)") //check user_id

	if err != nil {
		http.Error(w, "Failed to prepare SQL statement", http.StatusInternalServerError)
		return
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(stmt)
	_, err = stmt.Exec(title, body, status)
	if err != nil {
		http.Error(w, "Failed to execute SQL statement", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]string{"message": "Topic created successfully"}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("RegisterHandler: Error writing response after topic registered successfully : " + err.Error())
		return
	}
}

func deleteTopicHandler(w http.ResponseWriter, r *http.Request) {

	// Check method
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user_id, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Retrieve topicID from the request
	topicID := r.FormValue("topicID")

	var authorID int

	err := db.QueryRow("SELECT user_id FROM Topics_Table WHERE topic_id = ?", topicID).Scan(&authorID)
	if err != nil {
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}

	//verifier si IDuser est admin ou si c'est le proprietaire du topic
	isAdmin, err := checkAdminRights(user_id)

	if err != nil {
		http.Error(w, "Failed to check admin rights", http.StatusInternalServerError)
		return
	}
	if authorID != user_id {
		return
	} else if !isAdmin {
		return
	}

	// Parse data
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Delete topic from database
	stmt, err := db.Prepare("DELETE FROM Topics_Table WHERE id = ?")
	if err != nil {
		http.Error(w, "Failed to prepare SQL statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(topicID)
	if err != nil {
		http.Error(w, "Failed to execute SQL statement", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]string{"message": "Topic deleted successfully"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("Error writing response after topic deleted successfully: ", err)
		return
	}
}

func deleteCommentHandler(w http.ResponseWriter, r *http.Request) {

	// Check method
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user_id, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Retrieve topicID from the request
	topicID := r.FormValue("topicID")

	var authorID int

	err := db.QueryRow("SELECT user_id FROM Topics_Table WHERE topic_id = ?", topicID).Scan(&authorID)
	if err != nil {
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}

	//check if IDuser is admin or if it's the owner of the topic
	isAdmin, err := checkAdminRights(user_id)

	if err != nil {
		http.Error(w, "Failed to check admin rights", http.StatusInternalServerError)
		return
	}
	if authorID != user_id {
		return
	} else if !isAdmin {
		return
	}

	// Parse data
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Delete message from database
	stmt, err := db.Prepare("DELETE FROM Messages_Table WHERE id = ?")
	if err != nil {
		http.Error(w, "Failed to prepare SQL statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(topicID)
	if err != nil {
		http.Error(w, "Failed to execute SQL statement", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]string{"message": "Message deleted successfully"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("Error writing response after message deleted successfully: ", err)
		return
	}
}

func favoriteTopicHandler(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != "UPDATE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user_id, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Parse data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Retrieve topicID from the request
	topicID := r.FormValue("topicID")

	// Insert user into database
	stmt, err := db.Prepare("INSERT INTO follow  (topic_id, user_id) VALUES (?, ?)")
	if err != nil {
		http.Error(w, "Failed to prepare SQL statement", http.StatusInternalServerError)
		return
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(stmt) //
	_, err = stmt.Exec(topicID, user_id)
	if err != nil {
		http.Error(w, "Failed to execute SQL statement", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]string{"message": "Topic favorited successfully"}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("RegisterHandler: Error writing response after topic favorited successfully : " + err.Error())
		return
	}
}

func likeTopicHandler(w http.ResponseWriter, r *http.Request) {

	// Check method
	if r.Method != "UPDATE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user_id, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Parse data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Retrieve topicID from the request
	topicID := r.FormValue("topicID")

	// Check if the user has already liked the topic
	if alreadyLikedTopic(user_id, topicID) {
		err := neutralLikeComment(user_id, topicID) //bandite
		if err != nil {
			fmt.Printf("took of like on comment: %v\n", err)
			return
		}
		return
	}

	// Insert user into database
	stmt, err := db.Prepare("INSERT INTO react_topic (topic_id, user_id, status) VALUES (?, ?, 1)")
	if err != nil {
		http.Error(w, "Failed to prepare SQL statement", http.StatusInternalServerError)
		return
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(stmt) //
	_, err = stmt.Exec(topicID, user_id)
	if err != nil {
		http.Error(w, "Failed to execute SQL statement", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]string{"message": "Topic liked successfully"}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("RegisterHandler: Error writing response after topic liked successfully : " + err.Error())
		return
	}
}

func alreadyLikedTopic() bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM react_topic WHERE user_id = ? AND message_id = ? status = 1)", userID, messageID, status).Scan(&exists)
	if err != nil { //bandite
		fmt.Println(err)
		return false
	}
	return exists
}
