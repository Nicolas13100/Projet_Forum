package API

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go" // Import JWT library
	_ "github.com/go-sql-driver/mysql"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var (
	logged bool // Variable to track if user is logged in
)

// RegisterHandler Handler for confirming user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
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

	// Read the file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading avatar file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the image to a location on your server
	avatarPath := "/assets/images/userAvatar/"
	filename := username + filepath.Ext(r.FormValue("avatar"))
	err = os.WriteFile(avatarPath+filename, fileBytes, 0644)
	if err != nil {
		http.Error(w, "Error saving avatar file: "+err.Error(), http.StatusInternalServerError)
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
		response := map[string]string{"error": "Password incorrect, please use passwords with at least one lowercase letter, one uppercase letter, one digit, and a minimum length of 4 characters"}
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
	err = createUser(username, mail, hashedPassword, biography, avatarPath+filename)
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

// Secret key for JWT signing. It should be securely stored.
var jwtSecret = []byte("G45hUthd!3$gfdjHDfg@rT8p*3h$E%98")

// Claims Struct for JWT claims
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
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
		err = db.QueryRow("SELECT password FROM users_table WHERE username = ?", username).Scan(&storedPassword)
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
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
		next.ServeHTTP(w, r)
	}
}

// Handler for user logout
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	ResetUserValue()
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// Handler for dashboard
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if !logged {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	autoData := struct {
		Name   string
		Logged bool
	}{
		Name:   username,
		Logged: logged,
	}
	renderTemplate(w, "dashboard", autoData)
}

// Handler for gestion
func gestionHandler(w http.ResponseWriter, r *http.Request) {
	if !logged {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	data := struct {
		PlayerName string
		Logged     bool
	}{
		PlayerName: username,
		Logged:     logged,
	}
	renderTemplate(w, "gestion", data)
}

// Handler for changing login credentials
func changeLoginHandler(w http.ResponseWriter, r *http.Request) {
	if !logged {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	oldpassword := r.FormValue("oldpassword")
	newpassword := r.FormValue("newpassword")
	err := updateUserCredentials(username, oldpassword, newpassword)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Password updated successfully.")
	ResetUserValue()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Function to hash password
func hashPassword(password string) string {
	hasher := sha512.New()
	hasher.Write([]byte(password))
	hashedPassword := hasher.Sum(nil)
	return hex.EncodeToString(hashedPassword)
}

// Function to check if password matches hashed password
func checkPasswordHash(password, hash string) bool {
	hashedPassword := hashPassword(password)
	return hashedPassword == hash
}

// Function to load users from a file for register func
func loadUsersFromFile(filename string) error {
	// Check if the file exists
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// Create an empty users.json file if it doesn't exist
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println("error file.close login.go 225 ", err)
			}
		}(file)
	} else if err != nil {
		return err
	}

	// Check if the file is empty
	if fileInfo != nil && fileInfo.Size() == 0 {
		// File exists but is empty, so initialize users as an empty map
		users = make(map[string]User)
		return nil
	}

	// Load users from the file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("error file.close login.go 247 ", err)
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Check if the file contains valid JSON data
	if len(data) == 0 {
		// File is empty or doesn't contain valid JSON
		return nil
	}

	users = make(map[string]User)
	if err := json.Unmarshal(data, &users); err != nil {
		return err
	}

	return nil
}

// ResetUserValue Function to reset user values
func ResetUserValue() {
	logged = false
	username = ""
	password = ""
}

// Function to update user credentials
func updateUserCredentials(name, oldPassword, newPassword string) error {
	// Read the JSON file into memory
	raw, err := os.ReadFile("users.json")
	if err != nil {
		return err
	}

	// Define a struct that matches your JSON structure
	var data map[string]User // Map where keys are strings and values are User structs

	// Unmarshal the JSON into the defined struct
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	// Check if the user exists in the map
	user, exists := data[name]
	if !exists {
		return fmt.Errorf("user not found")
	}

	if !checkPasswordHash(oldPassword, user.Password) {
		return fmt.Errorf("incorrect password")
	}

	if newPassword != "" {
		// Update the password
		user.Password = hashPassword(newPassword)

		// Update the user in the map
		data[name] = user

		// Marshal the updated data back to JSON
		updatedJSON, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}

		// Write the updated JSON back to the file
		err = os.WriteFile("users.json", updatedJSON, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// Function to validate the password
func validatePassword(password string) bool {
	// Define the regex pattern
	pattern := `(?i)[a-z]+.*[A-Z]+.*\d+.+`

	// Compile the regex pattern into a regex object
	regex := regexp.MustCompile(pattern)

	// Check if the password matches the regex pattern
	return regex.MatchString(password)
}
