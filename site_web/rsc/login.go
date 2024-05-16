package API

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go" // Import JWT library
	_ "github.com/go-sql-driver/mysql"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

var (
	logged bool // Variable to track if user is logged in
)

// Secret key for JWT signing. It should be securely stored.
var jwtSecret = []byte("G45hUthd!3$gfdjHDfg@rT8p*3h$E%98")

// Claims Struct for JWT claims
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// RegisterHandler Handler for confirming user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data including file upload
	err := r.ParseMultipartForm(10 << 20) // Max size 10 MB
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Extract user registration data
	username := r.FormValue("username")
	password := r.FormValue("password")
	mail := r.FormValue("email")
	biography := r.FormValue("bio")
	file, _, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Error retrieving avatar: "+err.Error(), http.StatusBadRequest)
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
		http.Error(w, "Error saving profile picture: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if username or email already exist in the database
	var existingUser int
	err = db.QueryRow("SELECT COUNT(*) FROM users_table WHERE username = ? OR email = ?", username, mail).Scan(&existingUser)
	if err != nil {
		http.Error(w, "Error checking existing user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if existingUser > 0 {
		response := map[string]string{"error": "Username or email already exists"}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(jsonResponse)
		if err != nil {
			fmt.Println("RegisterHandler: Error writing response after checking existingUser > 0: " + err.Error())
			return
		}
		return
	}

	// Validate password
	isValid := validatePassword(password)
	if !isValid {
		response := map[string]string{"error": "Password incorrect, please use passwords with at least one lowercase letter, one uppercase letter, one digit, and a minimum length of 8 characters and a special character"}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(jsonResponse)
		if err != nil {
			fmt.Println("RegisterHandler: Error writing response after validatePassword : " + err.Error())
			return
		}
		return
	}

	// Hash password
	hashedPassword := hashPassword(password)

	// Create user in the database
	err = createUser(username, mail, hashedPassword, biography, profilePicPath+filename)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]string{"message": "User registered successfully"}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("RegisterHandler: Error writing response after User registered successfully : " + err.Error())
		return
	}
}

// Function to generate JWT token
func generateToken(userID int, username string) (string, error) {
	// Create the claims
	claims := &Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key and get the complete encoded token as a string
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// loginHandler Handler for user login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token   string `json:"token,omitempty"`
		Message string `json:"message,omitempty"`
		Error   string `json:"error,omitempty"`
	}

	// Check if the request method is POST
	if r.Method == "POST" {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
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
			// Handle error (e.g., user not found)
			err := json.NewEncoder(w).Encode(response{Error: "Invalid username or password"})
			if err != nil {
				fmt.Println("LoginHandler: Error writing response after login : " + err.Error())
				return
			}
			return
		}

		// Verify password
		if storedPassword != password {
			err := json.NewEncoder(w).Encode(response{Error: "Invalid username or password"})
			if err != nil {
				fmt.Println("LoginHandler: Error writing response after login invalid username or password: " + err.Error())
				return
			}
			return
		}

		// Authentication successful, generate token
		token, err := generateToken(userID, username)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Send the token in the response
		err = json.NewEncoder(w).Encode(response{Token: token, Message: "Authentication successful"})
		if err != nil {
			fmt.Println("LoginHandler: Error writing response after login Authentication successful: " + err.Error())
			return
		}
		return
	}

	// If request method is not POST
	err := json.NewEncoder(w).Encode(response{Error: "Method not allowed"})
	if err != nil {
		fmt.Println("Not a method of loginHandler : " + err.Error())
		return
	}
}

// Middleware function to authenticate requests
func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" && !blacklist.IsTokenBlacklisted(tokenString) {
			http.Error(w, "Unauthorized please log in", http.StatusUnauthorized)
			return
		}

		// Extract the token from the Authorization header
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify token validity
		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check expiration time
		claims, ok := token.Claims.(*Claims)
		if !ok || !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			http.Error(w, "Token has expired", http.StatusUnauthorized)
			return
		}

		// Token is valid and still within time
		// Pass userID to the next handler
		ctx := context.WithValue(r.Context(), "UserID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Handler for user logout
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Extract token from the request headers
	tokenString := r.Header.Get("Authorization")

	// Add the token to the blacklist
	blacklist.AddToken(tokenString, time.Now().Add(time.Hour*24)) // Blacklist token for 24H

	// Send JSON response
	response := map[string]string{
		"message": "Logged out successfully",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("Error writing JSON response:", err.Error())
	}
}

// Handler for dashboard
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("UserID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
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
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		return
	}

	// Send back the user data in the response
	jsonResponse, err := json.Marshal(userData)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("Error writing JSON response:", err.Error())
	}
}

// Handler for changing login credentials
func changeUserDataHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Parse form data including file upload
	err := r.ParseMultipartForm(10 << 20) // Max size 10 MB
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
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
		http.Error(w, "Error retrieving profile picture: "+err.Error(), http.StatusBadRequest)
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
		http.Error(w, "Error saving profile picture: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update the corresponding fields in the database
	_, err = db.Exec("UPDATE users_table SET username = ?, email = ?,password = ?, biography = ?,profile_pic = ? WHERE user_id = ?", newUsername, newEmail, newPassWord, newBiography, profilePicPath+filename, userID)
	if err != nil {
		http.Error(w, "Failed to update user data", http.StatusInternalServerError)
		return
	}

	// Send a success response
	response := map[string]string{
		"message": "User data updated successfully",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("Error writing JSON response:", err.Error())
	}
}
