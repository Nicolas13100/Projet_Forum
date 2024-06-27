package API

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

// RegisterHandler Handler for confirming user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data including file upload
	err := r.ParseMultipartForm(10 << 20) // Max size 10 MB
	if err != nil {
		fmt.Println(err)
		handleError(w, StatusBadRequest, "Failed to parse form data")
		return
	}
	// Extract user registration data
	username := r.FormValue("username")
	password := r.FormValue("password")
	mail := r.FormValue("mail")
	biography := r.FormValue("bio")

	// Validate username
	isUsernameValid := validateUsername(username)
	if !isUsernameValid {
		handleError(w, StatusBadRequest, "Username must be between 3 and 25 char long ")
		return
	}

	// Validate email
	isEmailValid := validateEmail(mail)
	if !isEmailValid {
		handleError(w, StatusBadRequest, "email requirements not met")
		return
	}

	var profilePicPath string
	var filename string

	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				fmt.Println("RegisterHandler: Error closing file : " + err.Error())
			}
		}(file)

		// Save the image to a location on your server
		profilePicPath = "/static/images/userAvatar/"
		filename = username + filepath.Ext(header.Filename)
		// To do, think if we save pics here or not
		//err = saveProfilePic(file, profilePicPath+filename)
		//if err != nil {
		//	fmt.Println(err)
		//	handleError(w, StatusInternalServerError, "Error saving profile picture")
		//	return
		//}
	}

	// Check if username or email already exist in the database
	var existingUser int
	err = db.QueryRow("SELECT COUNT(*) FROM users_table WHERE username = ? OR email = ?", username, mail).Scan(&existingUser)
	if err != nil {
		fmt.Println(err)
		handleError(w, StatusInternalServerError, "Error checking existing user")
		return
	}

	if existingUser > 0 {
		handleError(w, StatusBadRequest, "Username or email already exists")
		return
	}

	// Validate password
	isValid, message := validatePassword(password)
	if !isValid {
		handleError(w, StatusBadRequest, "Password requirements not met"+message)
		return
	}

	// Hash password
	hashedPassword := hashPassword(password)

	// Create user in the database
	err = createUser(username, mail, hashedPassword, biography, profilePicPath+filename)
	if err != nil {
		fmt.Println(err)
		handleError(w, StatusInternalServerError, "Error creating user")
		return
	}

	// Respond with success message
	response := APIResponse{Status: StatusOK, Message: "User registered successfully"}
	sendResponse(w, response)
}

// LoginHandler Handler for user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		handleError(w, StatusBadRequest, "Failed to parse form data")
		return
	}
	// Retrieve username and password from the form
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Hash password
	hashedPassword := hashPassword(password)

	// Query the database to check if the user exists and the password is correct
	var storedPassword string
	var userID int
	err = db.QueryRow("SELECT user_id, password FROM users_table WHERE username = ?", username).Scan(&userID, &storedPassword)
	if err != nil {
		fmt.Println(err)
		handleError(w, StatusUnauthorized, "Invalid username or password")
		return
	}
	// Verify password
	if storedPassword != hashedPassword {
		handleError(w, StatusUnauthorized, "Invalid username or password")
		return
	}

	// Authentication successful, generate token
	token, err := generateToken(userID)
	if err != nil {
		fmt.Println(err)
		handleError(w, StatusInternalServerError, "Failed to generate token")
		return
	}

	// Respond with success message and token
	response := APIResponse{Status: StatusOK, Message: "Authentication successful", Token: token}
	sendResponse(w, response)
	return
}

// LogoutHandler Handler for user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Extract token from the request headers
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		handleError(w, http.StatusUnauthorized, "Unauthorized: missing token")
		return
	}

	// Get the user ID from the context
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		handleError(w, http.StatusInternalServerError, "Failed to get user ID from context")
		return
	}

	// Delete the token from the database
	err := deleteTokenFromDB(userID, tokenString)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	// Respond with success message
	response := APIResponse{Status: http.StatusOK, Message: "Logged out successfully"}
	sendResponse(w, response)
}

// DashboardHandler Handler for dashboard
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		handleError(w, http.StatusInternalServerError, "User ID not found in context")
		return
	}

	// Query the database to fetch user data
	var userData struct {
		Username         string    `json:"username"`
		Email            string    `json:"email"`
		RegistrationDate time.Time `json:"registration_date"`
		LastLoginDate    time.Time `json:"last_login_date"`
		Biography        string    `json:"biography"`
		IsAdmin          bool      `json:"isAdmin"`
		IsModerator      bool      `json:"isModerator"`
		ProfilePic       string    `json:"profile_pic"`
	}

	err := db.QueryRow("SELECT username, email, registration_date, last_login_date, biography, isAdmin, isModerator, profile_pic FROM users_table WHERE user_id = ?", userID).Scan(
		&userData.Username,
		&userData.Email,
		&userData.RegistrationDate,
		&userData.LastLoginDate,
		&userData.Biography,
		&userData.IsAdmin,
		&userData.IsModerator,
		&userData.ProfilePic,
	)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to fetch user data")
		return
	}

	// Prepare JSON response
	jsonResponse, err := json.Marshal(userData)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to marshal JSON response")
		return
	}

	// Send back the user data in the response
	response := APIResponse{Status: http.StatusOK, Message: "Successfully fetched user data", JsonResp: jsonResponse}
	sendResponse(w, response)
}

// ChangeUserDataHandler Handler for changing login credentials
func ChangeUserDataHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		handleError(w, http.StatusInternalServerError, "User ID not found in context")
		return
	}

	// Parse form data including file upload
	err := r.ParseMultipartForm(10 << 20) // Max size 10 MB
	if err != nil {
		handleError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	// Retrieve new user data from the form
	newUsername := r.FormValue("username")
	newEmail := r.FormValue("email")
	newBiography := r.FormValue("biography")
	newPassword := hashPassword(r.FormValue("password")) // Assuming hashPassword function is defined

	// Handle profile picture upload
	file, _, err := r.FormFile("profile_pic")
	if err != nil {
		handleError(w, http.StatusBadRequest, "Error retrieving profile picture")
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err.Error())
		}
	}(file)

	// Save the profile picture to a location on your server
	profilePicPath := "/assets/images/userAvatar/"
	filename := strconv.Itoa(userID) + filepath.Ext(r.FormValue("profile_pic"))
	err = saveProfilePic(file, profilePicPath+filename)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Error saving profile picture")
		return
	}

	// Update the corresponding fields in the database
	_, err = db.Exec("UPDATE users_table SET username = ?, email = ?, password = ?, biography = ?, profile_pic = ? WHERE user_id = ?", newUsername, newEmail, newPassword, newBiography, profilePicPath+filename, userID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to update user data")
		return
	}

	response := APIResponse{Status: http.StatusOK, Message: "User data updated successfully"}
	sendResponse(w, response)
}
