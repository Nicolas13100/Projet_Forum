package API

import (
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

func Init() {
	InitBlacklist()
}

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
