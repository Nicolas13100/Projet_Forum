package API

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
		response := APIResponse{Status: http.StatusBadRequest, Message: "Failed to parse form data"}
		sendResponse(w, response)
		return
	}

	// Extract user registration data
	username := r.FormValue("username")
	password := r.FormValue("password")
	mail := r.FormValue("email")
	biography := r.FormValue("bio")
	file, _, err := r.FormFile("avatar")
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Error retrieving avatar"}
		sendResponse(w, response)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("RegisterHandler: Error closing file : " + err.Error())
		}
	}(file)

	// Save the image to a location on your server
	profilePicPath := "/assets/images/userAvatar/"
	filename := username + filepath.Ext(r.FormValue("avatar"))
	err = saveProfilePic(file, profilePicPath+filename)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Error saving profile picture"}
		sendResponse(w, response)
		return
	}

	// Check if username or email already exist in the database
	var existingUser int
	err = db.QueryRow("SELECT COUNT(*) FROM users_table WHERE username = ? OR email = ?", username, mail).Scan(&existingUser)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Error checking existing user"}
		sendResponse(w, response)
		return
	}

	if existingUser > 0 {
		response := APIResponse{Status: http.StatusBadRequest, Message: "Username or email already exists"}
		sendResponse(w, response)
		return
	}

	// Validate password
	isValid := validatePassword(password)
	if !isValid {
		response := APIResponse{Status: http.StatusBadRequest, Message: "Password incorrect, please use passwords with at least one lowercase letter, one uppercase letter, one digit, and a minimum length of 8 characters and a special character"}
		sendResponse(w, response)
		return
	}

	// Hash password
	hashedPassword := hashPassword(password)

	// Create user in the database
	err = createUser(username, mail, hashedPassword, biography, profilePicPath+filename)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Error creating user"}
		sendResponse(w, response)
		return
	}
	response := APIResponse{Status: http.StatusOK, Message: "User registered successfully"}
	sendResponse(w, response)
}

// LoginHandler Handler for user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	// Check if the request method is POST
	if r.Method == "POST" {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusBadRequest, Message: "Failed to parse form data"}
			sendResponse(w, response)
			return
		}

		// Retrieve username and password from the form
		username := r.FormValue("username")
		password := hashPassword(r.FormValue("password"))

		// Query the database to check if the user exists and the password is correct
		var storedPassword string
		var userID int
		err = db.QueryRow("SELECT user_id, password FROM users_table WHERE username = ?", username).Scan(&userID, &storedPassword)
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusTeapot, Message: "Invalid username or password"}
			sendResponse(w, response)
			return
		}

		// Verify password
		if storedPassword != password {
			response := APIResponse{Status: http.StatusTeapot, Message: "Invalid username or password"}
			sendResponse(w, response)
			return
		}

		// Authentication successful, generate token
		token, err := generateToken(userID)
		if err != nil {
			fmt.Println(err)
			response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to generate token"}
			sendResponse(w, response)
			return
		}
		response := APIResponse{Status: http.StatusOK, Message: "Authentication successful", Token: token}
		sendResponse(w, response)
		return
	}
	response := APIResponse{Status: http.StatusMethodNotAllowed, Message: "AMethod not allowed"}
	sendResponse(w, response)
}

// Authenticate Middleware function to Authenticate requests
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			response := APIResponse{Status: http.StatusUnauthorized, Message: "Unauthorized: please log in"}
			sendResponse(w, response)
			return
		}

		// Query the database to check token validity
		var response APIResponse
		var userID int
		var endDate time.Time
		err := db.QueryRow("SELECT user_id, end_date FROM tokens WHERE token = ?", tokenString).Scan(&userID, &endDate)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response = APIResponse{Status: http.StatusUnauthorized, Message: "Unauthorized: invalid token"}
			} else {
				response = APIResponse{Status: http.StatusInternalServerError, Message: "Internal server error"}
			}

			sendResponse(w, response)
			return
		}

		// Check if the token is expired
		if time.Now().After(endDate) {
			// Extend the token's expiration by 24 hours
			newEndDate := endDate.Add(24 * time.Hour)
			_, err = db.Exec("UPDATE tokens SET end_date = ? WHERE token = ?", newEndDate, tokenString)
			if err != nil {
				response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to extend token expiration, please re-log"}
				sendResponse(w, response)
				return
			}
			return
		}

		// Query the users_table to check if the user is an admin or a moderator
		var isAdmin, isModerator int
		err = db.QueryRow("SELECT isAdmin, isModerator FROM users_table WHERE user_id = ?", userID).Scan(&isAdmin, &isModerator)
		if err != nil {
			response = APIResponse{Status: http.StatusInternalServerError, Message: "Internal server error"}
			sendResponse(w, response)
			return
		}

		// Convert TINYINT values to booleans
		isAdminBool := isAdmin == 1
		isModeratorBool := isModerator == 1

		// Token is valid and not expired
		// Pass userID, isAdminBool, and isModeratorBool to the next handler
		ctx := context.WithValue(r.Context(), "UserID", userID)
		ctx = context.WithValue(ctx, "IsAdmin", isAdminBool)
		ctx = context.WithValue(ctx, "IsModerator", isModeratorBool)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// LogoutHandler Handler for user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Extract token from the request headers
	tokenString := r.Header.Get("Authorization")
	// Get the user ID from the context
	userID := r.Context().Value("UserID").(int)

	err := deleteTokenFromDB(userID, tokenString)
	if err != nil {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to logout"}
		sendResponse(w, response)
		return
	}

	response := APIResponse{Status: http.StatusOK, Message: "Logged out successfully"}
	sendResponse(w, response)
}

// DashboardHandler Handler for dashboard
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("UserID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
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
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to fetch user data"}
		sendResponse(w, response)
		return
	}

	// Send back the user data in the response
	jsonResponse, err := json.Marshal(userData)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to marshal JSON response"}
		sendResponse(w, response)
		return
	}
	response := APIResponse{Status: http.StatusOK, Message: "Correctly fetched user data", JsonResp: jsonResponse}
	sendResponse(w, response)
}

// ChangeUserDataHandler Handler for changing login credentials
func ChangeUserDataHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		response := APIResponse{Status: http.StatusInternalServerError, Message: "User ID not found in context"}
		sendResponse(w, response)
		return
	}

	// Parse form data including file upload
	err := r.ParseMultipartForm(10 << 20) // Max size 10 MB
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Failed to parse form data"}
		sendResponse(w, response)
		return
	}

	// Retrieve new user data from the form
	newUsername := r.FormValue("username")
	newEmail := r.FormValue("email")
	newBiography := r.FormValue("biography")
	newPassWord := hashPassword(r.FormValue("password"))
	// You can retrieve other fields similarly

	// Handle profile picture upload
	file, _, err := r.FormFile("profile_pic")
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusBadRequest, Message: "Error retrieving profile picture"}
		sendResponse(w, response)
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
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Error saving profile picture"}
		sendResponse(w, response)
		return
	}

	// Update the corresponding fields in the database
	_, err = db.Exec("UPDATE users_table SET username = ?, email = ?,password = ?, biography = ?,profile_pic = ? WHERE user_id = ?", newUsername, newEmail, newPassWord, newBiography, profilePicPath+filename, userID)
	if err != nil {
		fmt.Println(err)
		response := APIResponse{Status: http.StatusInternalServerError, Message: "Failed to update user data"}
		sendResponse(w, response)
		return
	}
	response := APIResponse{Status: http.StatusOK, Message: "User data updated successfully"}
	sendResponse(w, response)
}
