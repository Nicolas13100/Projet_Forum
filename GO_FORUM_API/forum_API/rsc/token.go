package API

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func GetUserIdByToken(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		sendResponse(w, APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	token := vars["token"]

	if token == "" {
		sendResponse(w, APIResponse{Status: http.StatusBadRequest, Message: "Invalid token"})
		return
	}

	var userID int
	err := db.QueryRow(`
        SELECT user_id
        FROM tokens
        WHERE token = ?`, token).Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendResponse(w, APIResponse{Status: http.StatusNotFound, Message: "Topic not found"})
		} else {
			log.Printf("Error querying database: %v", err)
			sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		}
		return
	}

	sendResponse(w, APIResponse{
		Status:  http.StatusOK,
		Message: "Success",
		UserID:  userID,
	})
}

func DeleteToken(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		sendResponse(w, APIResponse{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	token := vars["token"]

	if token == "" {
		sendResponse(w, APIResponse{Status: http.StatusBadRequest, Message: "Invalid token"})
		return
	}

	// Execute DELETE query
	_, err := db.Exec(`
		DELETE FROM tokens
		WHERE token = ?`, token)

	if err != nil {
		sendResponse(w, APIResponse{Status: http.StatusInternalServerError, Message: "Failed to delete token"})
		return
	}

	sendResponse(w, APIResponse{
		Status:  http.StatusOK,
		Message: "Token deleted successfully",
	})
}
