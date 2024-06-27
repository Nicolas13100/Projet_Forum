package API

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"strconv"
)

func CreateTopicHandler(w http.ResponseWriter, r *http.Request) {

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

func DeleteTopicHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {

	// Check method
	if r.Method != "DELETE" {
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
	isAdmin, err := checkAdminRights(userId)

	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to check admin rights"}
		sendResponse(w, response)
		return
	}
	if authorID != userId || !isAdmin {
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

func FavoriteTopicHandler(w http.ResponseWriter, r *http.Request) {
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

func GetAllTopic(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		response := APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"}
		sendResponse(w, response)
		return
	}

	// Extract page and pageSize from request URL
	vars := mux.Vars(r)
	pageStr := vars["page"]
	pageSizeStr := vars["pageSize"]

	// Parse page and pageSize from string to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Default to page 1 if page parameter is invalid
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10 // Default page size to 10 if pageSize parameter is invalid
	}
	offset := (page - 1) * pageSize

	// Query the topics with pagination
	rows, err := db.Query("SELECT topic_id, title, body, creation_date, status, is_private, user_id FROM Topics_Table LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		log.Println("Error querying database:", err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
		sendResponse(w, response)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("Error closing rows:", err)
		}
	}()

	var topics []Topic
	for rows.Next() {
		var topic Topic
		var isPrivate int
		if err := rows.Scan(&topic.TopicID, &topic.Title, &topic.Body, &topic.CreationDate, &topic.Status, &isPrivate, &topic.UserID); err != nil {
			log.Println("Error scanning rows:", err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
			sendResponse(w, response)
			return
		}
		topic.IsPrivate = isPrivate == 1
		topics = append(topics, topic)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error with rows iteration:", err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
		sendResponse(w, response)
		return
	}

	// Count total number of rows (topics) in the table
	var totalCount int
	err = db.QueryRow("SELECT COUNT(*) FROM Topics_Table").Scan(&totalCount)
	if err != nil {
		log.Println("Error counting total rows:", err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
		sendResponse(w, response)
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	// Prepare response
	responseData := map[string]interface{}{
		"status":      http.StatusOK,
		"message":     "Success",
		"topics":      topics,
		"totalCount":  totalCount,
		"totalPages":  totalPages,
		"currentPage": page,
	}

	response := APIResponse{Status: http.StatusOK, Message: "Success", Resp: responseData}
	sendResponse(w, response)
}

func GetTopic(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		response := APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"}
		sendResponse(w, response)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	// Query the database for the topic with the given ID
	topic := Topic{}
	err := db.QueryRow("SELECT topic_id, title, body, creation_date, status, is_private, user_id FROM Topics_Table WHERE topic_id = ?", id).
		Scan(&topic.TopicID, &topic.Title, &topic.Body, &topic.CreationDate, &topic.Status, &topic.IsPrivate, &topic.UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no topic found with the given ID, return 404 Not Found
			http.Error(w, "Topic not found", http.StatusNotFound)
		} else {
			// For other errors, return 500 Internal Server Error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Convert topic to JSON
	topicJSON, err := json.Marshal(topic)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	response := APIResponse{Status: http.StatusOK, Message: "Success", JsonResp: topicJSON}
	sendResponse(w, response)
}

func GetTopicTagsNamesByTopicId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"}
		sendResponse(w, response)
		return
	}

	// Extract topicId from URL parameter
	vars := mux.Vars(r)
	topicIdStr := vars["id"]

	topicId, err := strconv.Atoi(topicIdStr)
	if err != nil || topicId < 1 {
		log.Println("Invalid topic ID:", topicIdStr)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Invalid topic ID"}
		sendResponse(w, response)
		return
	}

	rows, err := db.Query("SELECT t.tag_id, t.tag_name FROM Tags_Table t JOIN TopicTags tt ON t.tag_id = tt.tag_id WHERE tt.topic_id = ?;", topicId)

	if err != nil {
		log.Println("Error querying database:", err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
		sendResponse(w, response)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("Error closing rows:", err)
		}
	}()

	var tagInfos []TagInfo

	// Process rows
	for rows.Next() {
		var tagID int
		var tagName string
		if err := rows.Scan(&tagID, &tagName); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		tagInfo := TagInfo{
			TagID:   tagID,
			TagName: tagName,
		}
		tagInfos = append(tagInfos, tagInfo)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
		sendResponse(w, response)
		return
	}

	// Prepare response with all tag names
	response := APIResponse{Status: http.StatusOK, Data: tagInfos}
	sendResponse(w, response)
}
