package API

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func commentHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	// Parse form values
	topicID, err := strconv.Atoi(r.FormValue("topicID"))
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Invalid topic ID"}
		sendResponse(w, response)
		return
	}

	commentID, err := strconv.Atoi(r.FormValue("commentID"))
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Invalid comment ID"}
		sendResponse(w, response)
		return
	}

	comment := r.FormValue("comment")

	// Perform actions based on the provided parameters
	if topicID != 0 {
		err := createTopicComment(userID, topicID, comment)
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to create topic comment"}
			sendResponse(w, response)
			return
		}
		response := APIResponse{Status: http.StatusOK, Message: "Comment on topic created successfully"}
		sendResponse(w, response)
		return
	}

	if commentID != 0 {
		// Find the topic ID of the base message
		baseTopicID, err := getBaseMessageTopicID(commentID)
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to get topic ID of base message"}
			sendResponse(w, response)
			return
		}
		err = createReplyToComment(userID, commentID, baseTopicID, comment)
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "couldn't creat reply to comment"}
			sendResponse(w, response)
			return
		}
		response := APIResponse{Status: http.StatusOK, Message: "Reply to comment created successfully"}
		sendResponse(w, response)
		return
	}

	response := APIResponse{Status: http.StatusBadRequest, Message: "Invalid request"}
	sendResponse(w, response)
}

func sendResponse(w http.ResponseWriter, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("send response error:", err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "API send response error"}
		sendResponse(w, response)
		return
	}
}

func createTopicComment(userID int, topicID int, comment string) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO Messages_Table (body, topic_id, user_id) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println("Error closing stmt:", err)
		}
	}(stmt)

	// Execute the SQL statement
	_, err = stmt.Exec(comment, topicID, userID)
	if err != nil {
		return err
	}

	return nil
}

func createReplyToComment(userID int, commentID int, topicID int, comment string) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO Messages_Table (body, base_message_id, user_id, topic_id) VALUES (?, ?, ?,?)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println("Error closing stmt:", err)
		}
	}(stmt)

	// Execute the SQL statement
	_, err = stmt.Exec(comment, commentID, userID, topicID)
	if err != nil {
		return err
	}

	return nil
}

func getBaseMessageTopicID(commentID int) (int, error) {
	var topicID int
	err := db.QueryRow("SELECT topic_id FROM Messages_Table WHERE message_id = ?", commentID).Scan(&topicID)
	if err != nil {
		return 0, err
	}
	return topicID, nil
}

func updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}
	messageID, err := strconv.Atoi(r.FormValue("messageID"))
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Invalid message ID"}
		sendResponse(w, response)
		return
	}
	admin, err := checkAdminRights(userID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to check admin rights"}
		sendResponse(w, response)
		return
	}
	// Query the database to retrieve the authorID associated with the topicID
	var authorID int
	err = db.QueryRow("SELECT user_id FROM Messages_Table WHERE message_id = ?", messageID).Scan(&authorID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Impossible to see if user id author"}
		sendResponse(w, response)
		return
	}

	if authorID != userID || !admin {
		response := APIResponse{Status: http.StatusForbidden, Message: "You do not have access to this topic"}
		sendResponse(w, response)
		return
	}

	// Update the comment in the database
	_, err = db.Exec("UPDATE Messages_Table SET body = ? WHERE message_id = ?", r.FormValue("body"), messageID)
	if err != nil {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to update comment"}
		sendResponse(w, response)
		return
	}

	response := APIResponse{Status: http.StatusOK, Message: "Comment updated successfully"}
	sendResponse(w, response)
}
