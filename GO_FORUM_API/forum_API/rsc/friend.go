package API

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func AddFriendHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	askedIDStr := r.FormValue("askedID")
	if askedIDStr == "" {
		response := APIResponse{Status: http.StatusBadRequest, Message: "askedID parameter is required"}
		sendResponse(w, response)
		return
	}

	askedID, err := strconv.Atoi(askedIDStr)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Invalid askedID parameter"}
		sendResponse(w, response)
		return
	}

	// Check if the friendship request already exists
	// You need to replace db with your actual database connection
	var status int
	err = db.QueryRow("SELECT status FROM friendship WHERE sender_id = ? AND reciver_id = ?", userID, askedID).Scan(&status)
	if err == nil {
		switch status {
		case 0:
			response := APIResponse{Status: http.StatusBadRequest, Message: "Friend request already sent"}
			sendResponse(w, response)
		case 1:
			response := APIResponse{Status: http.StatusBadRequest, Message: "Friend request already accepted"}
			sendResponse(w, response)
		case 2:
			response := APIResponse{Status: http.StatusBadRequest, Message: "Friend request already refused"}
			sendResponse(w, response)
		}
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Error checking friendship request"}
		sendResponse(w, response)
		return
	}

	// Insert a new friendship request
	createdAt := time.Now()
	_, err = db.Exec("INSERT INTO friendship (sender_id, reciver_id, status, created_at, updated_at) VALUES (?, ?, 0, ?, ?)", userID, askedID, createdAt, createdAt)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Error inserting friendship request"}
		sendResponse(w, response)
		return
	}
	response := APIResponse{Status: http.StatusCreated, Message: "Friendship request successfully sent"}
	sendResponse(w, response)
}

func AcceptFriendHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	askedIDStr := r.FormValue("askedID")
	if askedIDStr == "" {
		response := APIResponse{Status: http.StatusBadRequest, Message: "askedID parameter is required"}
		sendResponse(w, response)
		return
	}

	askedID, err := strconv.Atoi(askedIDStr)
	if err != nil {
		response := APIResponse{Status: http.StatusBadRequest, Message: "Invalid askedID parameter"}
		sendResponse(w, response)
		return
	}

	// Check if there is a pending friendship request from the user represented by askedID
	var status int
	err = db.QueryRow("SELECT status FROM friendship WHERE sender_id = ? AND reciver_id = ? AND status = 0", askedID, userID).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "No pending friend request found"}
		sendResponse(w, response)
		return
	} else if err != nil {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Error checking friendship request"}
		sendResponse(w, response)
		return
	}

	// Update the status of the friendship request to accepted
	updatedAt := time.Now()
	_, err = db.Exec("UPDATE friendship SET status = 1, updated_at = ? WHERE sender_id = ? AND reciver_id = ?", updatedAt, askedID, userID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Error accepting friendship request"}
		sendResponse(w, response)
		return
	}
	response := APIResponse{Status: http.StatusOK, Message: "Friend request accepted successfully"}
	sendResponse(w, response)
}

func DeclineFriendHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	askedIDStr := r.FormValue("askedID")
	if askedIDStr == "" {
		response := APIResponse{Status: http.StatusBadRequest, Message: "askedID parameter is required"}
		sendResponse(w, response)
		return
	}

	askedID, err := strconv.Atoi(askedIDStr)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Invalid askedID parameter"}
		sendResponse(w, response)
		return
	}

	// Check if there is a pending friendship request from the user represented by askedID
	var status int
	err = db.QueryRow("SELECT status FROM friendship WHERE sender_id = ? AND reciver_id = ? AND status = 0", askedID, userID).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "No pending friend request found"}
		sendResponse(w, response)
		return
	} else if err != nil {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Error checking friendship request"}
		sendResponse(w, response)
		return
	}

	// Update the status of the friendship request to refused
	updatedAt := time.Now()
	_, err = db.Exec("UPDATE friendship SET status = 2, updated_at = ? WHERE sender_id = ? AND reciver_id = ?", updatedAt, askedID, userID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Error declining friendship request"}
		sendResponse(w, response)
		return
	}
	response := APIResponse{Status: http.StatusOK, Message: "Friend request declined successfully"}
	sendResponse(w, response)
}

func DeleteFriendHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	friendIDStr := r.FormValue("friendID")
	if friendIDStr == "" {
		response := APIResponse{Status: http.StatusBadRequest, Message: "friendID parameter is required"}
		sendResponse(w, response)
		return
	}

	friendID, err := strconv.Atoi(friendIDStr)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Invalid friendID parameter"}
		sendResponse(w, response)
		return
	}

	// Check if the friendship exists
	var status int
	err = db.QueryRow("SELECT status FROM friendship WHERE (sender_id = ? AND reciver_id = ?) OR (sender_id = ? AND reciver_id = ?)", userID, friendID, friendID, userID).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "No friendship found with the specified user"}
		sendResponse(w, response)
		return
	} else if err != nil {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Error checking friendship"}
		sendResponse(w, response)
		return
	}

	// Delete the friendship
	_, err = db.Exec("DELETE FROM friendship WHERE (sender_id = ? AND reciver_id = ?) OR (sender_id = ? AND reciver_id = ?)", userID, friendID, friendID, userID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Error deleting friendship"}
		sendResponse(w, response)
		return
	}
	response := APIResponse{Status: http.StatusOK, Message: "Friend successfully deleted"}
	sendResponse(w, response)
}
