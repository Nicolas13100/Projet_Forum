package API

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func addFriendHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	askedIDStr := r.FormValue("askedID")
	if askedIDStr == "" {
		http.Error(w, "askedID parameter is required", http.StatusBadRequest)
		return
	}

	askedID, err := strconv.Atoi(askedIDStr)
	if err != nil {
		http.Error(w, "Invalid askedID parameter", http.StatusBadRequest)
		return
	}

	// Check if the friendship request already exists
	// You need to replace db with your actual database connection
	var status int
	err = db.QueryRow("SELECT status FROM friendship WHERE sender_id = ? AND reciver_id = ?", userID, askedID).Scan(&status)
	if err == nil {
		switch status {
		case 0:
			http.Error(w, "Friend request already sent", http.StatusBadRequest)
		case 1:
			http.Error(w, "Friend request already accepted", http.StatusBadRequest)
		case 2:
			http.Error(w, "Friend request already refused", http.StatusBadRequest)
		}
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Error checking friendship request", http.StatusInternalServerError)
		return
	}

	// Insert a new friendship request
	createdAt := time.Now()
	_, err = db.Exec("INSERT INTO friendship (sender_id, reciver_id, status, created_at, updated_at) VALUES (?, ?, 0, ?, ?)", userID, askedID, createdAt, createdAt)
	if err != nil {
		http.Error(w, "Error inserting friendship request", http.StatusInternalServerError)
		return
	}

	// Friendship request successfully sent
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("Friend request sent successfully"))
	if err != nil {
		fmt.Println("friendship write", err)
		return
	}
}

func acceptFriendHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	askedIDStr := r.FormValue("askedID")
	if askedIDStr == "" {
		http.Error(w, "askedID parameter is required", http.StatusBadRequest)
		return
	}

	askedID, err := strconv.Atoi(askedIDStr)
	if err != nil {
		http.Error(w, "Invalid askedID parameter", http.StatusBadRequest)
		return
	}

	// Check if there is a pending friendship request from the user represented by askedID
	var status int
	err = db.QueryRow("SELECT status FROM friendship WHERE sender_id = ? AND reciver_id = ? AND status = 0", askedID, userID).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "No pending friend request found", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Error checking friendship request", http.StatusInternalServerError)
		return
	}

	// Update the status of the friendship request to accepted
	updatedAt := time.Now()
	_, err = db.Exec("UPDATE friendship SET status = 1, updated_at = ? WHERE sender_id = ? AND reciver_id = ?", updatedAt, askedID, userID)
	if err != nil {
		http.Error(w, "Error accepting friendship request", http.StatusInternalServerError)
		return
	}

	// Friendship request successfully accepted
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Friend request accepted successfully"))
	if err != nil {
		fmt.Println("accept friendship write", err)
		return
	}
}

func declineFriendHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	askedIDStr := r.FormValue("askedID")
	if askedIDStr == "" {
		http.Error(w, "askedID parameter is required", http.StatusBadRequest)
		return
	}

	askedID, err := strconv.Atoi(askedIDStr)
	if err != nil {
		http.Error(w, "Invalid askedID parameter", http.StatusBadRequest)
		return
	}

	// Check if there is a pending friendship request from the user represented by askedID
	var status int
	err = db.QueryRow("SELECT status FROM friendship WHERE sender_id = ? AND reciver_id = ? AND status = 0", askedID, userID).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "No pending friend request found", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Error checking friendship request", http.StatusInternalServerError)
		return
	}

	// Update the status of the friendship request to refused
	updatedAt := time.Now()
	_, err = db.Exec("UPDATE friendship SET status = 2, updated_at = ? WHERE sender_id = ? AND reciver_id = ?", updatedAt, askedID, userID)
	if err != nil {
		http.Error(w, "Error declining friendship request", http.StatusInternalServerError)
		return
	}

	// Friendship request successfully declined
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Friend request declined successfully"))
	if err != nil {
		fmt.Println("decline friendship write", err)
		return
	}
}

func deleteFriendHandler(w http.ResponseWriter, r *http.Request) {
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
		response := APIResponse{Status: http.StatusBadRequest, Message: "Invalid friendID parameter"}
		sendResponse(w, response)
		return
	}

	// Check if the friendship exists
	var status int
	err = db.QueryRow("SELECT status FROM friendship WHERE (sender_id = ? AND reciver_id = ?) OR (sender_id = ? AND reciver_id = ?)", userID, friendID, friendID, userID).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
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
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Error deleting friendship"}
		sendResponse(w, response)
		return
	}

	// Friendship successfully deleted
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Friend successfully deleted"))
	if err != nil {
		fmt.Println("delete friendship write", err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Error writing ok answer"}
		sendResponse(w, response)
		return
	}
}
