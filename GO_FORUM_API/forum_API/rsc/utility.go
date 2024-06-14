package API

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// ////
func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	// Taken from hangman
	tmpl, err := template.New(tmplName).Funcs(template.FuncMap{"join": join, "contains": containsString}).ParseFiles("forum_API/Template/" + tmplName + ".html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func join(s []string, sep string) string {
	// same
	return strings.Join(s, sep)
}
func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

///////

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
func validatePassword(password string) bool {
	// Define the regex pattern
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`

	// Compile the regex pattern into a regex object
	regex := regexp.MustCompile(pattern)

	// Check if the password matches the regex pattern
	return regex.MatchString(password)
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
