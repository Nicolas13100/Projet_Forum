package API

import (
	"database/sql"
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

	if title == "" || body == "" || status == "" || userID != 0 {
		response := APIResponse{Status: http.StatusBadRequest, Message: "Missing required form fields"}
		sendResponse(w, response)
		return
	}

	// Insert user into database
	stmt, err := db.Prepare("INSERT INTO Topics_Table (title, body, status, user_id) VALUES (?, ?, ?,?)") //check user_id

	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to prepare SQL statement"}
		sendResponse(w, response)
		return
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to close SQL statement"}
			sendResponse(w, response)
		}
	}(stmt)
	_, err = stmt.Exec(title, body, status, userID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to execute SQL statement"}
		sendResponse(w, response)
		return
	}
	response := APIResponse{Status: http.StatusOK, Message: "Topic created successfully"}
	sendResponse(w, response)
}

func deleteTopicHandler(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != "DELETE" {
		response := APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"}
		sendResponse(w, response)
		return
	}

	// Parse data
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Failed to parse form data"}
		sendResponse(w, response)
		return
	}

	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	// Retrieve topicID from the request
	topicID := r.FormValue("topicID")

	var authorID int

	err = db.QueryRow("SELECT user_id FROM Topics_Table WHERE topic_id=?", topicID).Scan(&authorID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusNotFound, Message: "Topic not found"}
		sendResponse(w, response)
		return
	}

	// Check authorization
	var isAdmin bool
	isAdmin, err = checkAdminRights(userId)

	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to check admin rights"}
		sendResponse(w, response)
		return
	}
	if authorID != userId && !isAdmin {
		response := APIResponse{Status: http.StatusUnauthorized, Message: "No rights to delete topic"}
		sendResponse(w, response)
		return
	}

	// Delete topic from database
	stmt, err := db.Prepare("DELETE FROM Topics_Table WHERE topic_id = ?")
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to prepare SQL statement"}
		sendResponse(w, response)
		return
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to close SQL statement"}
			sendResponse(w, response)
		}
	}(stmt)

	_, err = stmt.Exec(topicID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to execute SQL statement"}
		sendResponse(w, response)
		return
	}

	// Return success response
	response := APIResponse{Status: http.StatusOK, Message: "Topic deleted successfully"}
	sendResponse(w, response)
}

func deleteCommentHandler(w http.ResponseWriter, r *http.Request) {

	// Check method
	if r.Method != "DELETE" {
		response := APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"}
		sendResponse(w, response)
		return
	}

	user_id, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	// Retrieve topicID from the request
	topicID := r.FormValue("topicID")

	var authorID int
	err := db.QueryRow("SELECT user_id FROM Topics_Table WHERE topic_id=?", topicID).Scan(&authorID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusNotFound, Message: "Topic not found"}
		sendResponse(w, response)
		return
	}

	//check if userID is admin or if it's the owner of the topic
	isAdmin, err := checkAdminRights(user_id)

	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to check admin rights"}
		sendResponse(w, response)
		return
	}
	if authorID != user_id || !isAdmin {
		response := APIResponse{Status: http.StatusUnauthorized, Message: "Not allowed to delete comment"}
		sendResponse(w, response)
		return
	}

	// Parse data
	err = r.ParseForm()
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to parse form data in delete comment"}
		sendResponse(w, response)
		return
	}

	// Delete message from database
	stmt, err := db.Prepare("DELETE FROM messages_table WHERE topic_id = ?")
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to prepare SQL statement"}
		sendResponse(w, response)
		return
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to close SQL statement"}
			sendResponse(w, response)
		}
	}(stmt)

	_, err = stmt.Exec(topicID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to execute SQL statement"}
		sendResponse(w, response)
		return
	}

	// Return success response
	response := APIResponse{Status: http.StatusOK, Message: "Message deleted successfully"}
	sendResponse(w, response)
}

func favoriteTopicHandler(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != "UPDATE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		response := APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"}
		sendResponse(w, response)
		return
	}

	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	// Parse data
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Failed to parse form data"}
		sendResponse(w, response)
		return
	}

	// Retrieve topicID from the request
	topicID := r.FormValue("topicID")

	// Insert user into database
	stmt, err := db.Prepare("INSERT INTO follow (user_id, topic_id) VALUE (?,?)")
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to prepare SQL statement"}
		sendResponse(w, response)
		return
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to close SQL statement"}
			sendResponse(w, response)
		}
	}(stmt) //
	_, err = stmt.Exec(userId, topicID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to execute SQL statement"}
		sendResponse(w, response)
		return
	}

	// Return success response
	response := APIResponse{Status: http.StatusOK, Message: "Topic added to favorite successfully"}
	sendResponse(w, response)
}
