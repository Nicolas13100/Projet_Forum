package API

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"time"
	"unicode"
)

const (
	StatusOK                  = http.StatusOK
	StatusBadRequest          = http.StatusBadRequest
	StatusUnauthorized        = http.StatusUnauthorized
	StatusInternalServerError = http.StatusInternalServerError
	StatusMethodNotAllowed    = http.StatusMethodNotAllowed
)

// contextKey defines a custom type for context keys
type contextKey string

const (
	userIDKey      contextKey = "UserID"
	isAdminKey     contextKey = "IsAdmin"
	isModeratorKey contextKey = "IsModerator"
)

// Function to save profile picture to the server
func saveProfilePic(file multipart.File, path string) error {
	// Create the profile picture file
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(dst)

	// Copy the file content
	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	return nil
}

// Function to hash password
func hashPassword(password string) string {
	hasher := sha512.New()
	hasher.Write([]byte(password))
	hashedPassword := hasher.Sum(nil)
	return hex.EncodeToString(hashedPassword)
}

// Function to validate the password
func validatePassword(password string) (bool, string) {
	// Check length
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}

	var (
		hasUpperCase bool
		hasLowerCase bool
		hasDigit     bool
		hasSpecial   bool
	)

	// Check each character in the password
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpperCase = true
		case unicode.IsLower(char):
			hasLowerCase = true
		case unicode.IsDigit(char):
			hasDigit = true
		case isSpecialChar(char):
			hasSpecial = true
		}
	}

	// Check if all criteria are met
	if !hasUpperCase {
		return false, "Password must contain at least one uppercase letter"
	}
	if !hasLowerCase {
		return false, "Password must contain at least one lowercase letter"
	}
	if !hasDigit {
		return false, "Password must contain at least one digit"
	}
	if !hasSpecial {
		return false, "Password must contain at least one special character"
	}

	// If all checks pass
	return true, "Password is valid"
}

// Function to check if a character is one of the allowed special characters
func isSpecialChar(char rune) bool {
	switch char {
	case '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '_', '=', '+', '~',
		'[', ']', '{', '}', '|', ';', ':', '\'', '"', ',', '.', '<', '>', '/', '?':
		return true
	default:
		return false
	}
}

// Function to validate the password
func validateUsername(username string) bool {
	return len(username) >= 3 && len(username) <= 25
}

// Function to validate the password
func validateEmail(mail string) bool {
	// Define the regex pattern
	pattern := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`

	// Compile the regex pattern into a regex object
	regex := regexp.MustCompile(pattern)

	// Step-by-step validation
	if len(mail) == 0 {
		fmt.Println("Error: Email address is empty")
		return false
	}

	// Check if the email contains exactly one @ symbol
	atCount := 0
	for _, char := range mail {
		if char == '@' {
			atCount++
		}
	}
	if atCount != 1 {
		fmt.Println("Error: Email address must contain exactly one '@' symbol")
		return false
	}

	// Split the email into local and domain parts
	parts := regexp.MustCompile(`@`).Split(mail, 2)
	localPart, domainPart := parts[0], parts[1]

	// Validate the local part
	localPattern := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+$`
	if !regexp.MustCompile(localPattern).MatchString(localPart) {
		fmt.Println("Error: Invalid characters in the local part")
		return false
	}

	// Validate the domain part
	domainPattern := `^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	if !regexp.MustCompile(domainPattern).MatchString(domainPart) {
		fmt.Println("Error: Invalid domain format")
		return false
	}

	// Final regex match for overall email format
	if !regex.MatchString(mail) {
		fmt.Println("Error: Email does not match the overall pattern")
		return false
	}
	return true
}

// deleteTokenFromDB deletes the token for the specified user ID
func deleteTokenFromDB(userID int, token string) error {
	// Assuming you have a DB connection set up as `db`
	query := "DELETE FROM tokens WHERE user_id = ? AND token = ?"
	_, err := db.Exec(query, userID, token)
	return err
}

// Function to generate JWT token
func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

func generateToken(userID int) (string, error) {
	// Generate a random 86-character token
	token, err := generateRandomToken(86)
	if err != nil {
		return "", err
	}

	// Define the token expiration time (24 hours from now)
	endDate := time.Now().Add(24 * time.Hour)

	// Store the token in the database
	query := `INSERT INTO tokens (user_id, end_date, token) VALUES (?, ?, ?)`
	_, err = db.Exec(query, userID, endDate, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func handleError(w http.ResponseWriter, status int, message string) {
	response := APIResponse{Status: status, Message: message}
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
