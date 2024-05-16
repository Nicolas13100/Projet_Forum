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
		response := APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"}
		sendResponse(w, response)
		return
	}

	//parse data
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		response := APIResponse{Status: http.StatusBadRequest, Message: "Failed to parse form data"}
		sendResponse(w, response)
		return
	}

	//check if all the data are here
	title := r.FormValue("title")
	body := r.FormValue("body")
	status := r.FormValue("status")
	//check image
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	if title == "" || body == "" || status == "" || userID != 0 )
		response := APIResponse{Status: http.StatusBadRequest, Message: "Missing required form fields"}
		sendResponse(w, response)
		return
	}

	// Insert user into database
	stmt, err := db.Prepare("INSERT INTO Topics_Table (title, body, status, user_id) VALUES (?, ?, ?,?)") //check user_id

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
	_, err = stmt.Exec(title, body, status, userID)
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
