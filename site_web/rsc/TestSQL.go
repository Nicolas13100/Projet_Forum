package API

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/liveuser")
	if err != nil {
		return err
	}
	// Test the connection
	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}

var store *sessions.CookieStore

func initSessionStore() {
	// Replace "your-secret-key" with a strong, random key
	store = sessions.NewCookieStore([]byte("your-secret-key"))
	store.Options.MaxAge = 60 * 60 // One hour session duration (customize)
}

func createUser(username, password string) error {
	hashedPassword := hashPassword(password)
	// Insert user into database
	stmt, err := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, hashedPassword)
	return err
}

func getUserByUsername(username string) (User, error) {
	var user User
	// Retrieve user from database
	stmt, err := db.Prepare("SELECT id, username, password FROM users WHERE username = ?")
	if err != nil {
		return user, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return user, fmt.Errorf("user not found")
	}
	return user, err
}
func createSession(w http.ResponseWriter, r *http.Request, userID int64) error {
	session, err := store.Get(r, "user-session")
	if err != nil {
		return err
	}
	session.Values["user_id"] = userID
	err = session.Save(r, w)
	return err
}

func getSession(r *http.Request) (int64, error) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		return 0, err
	}
	userID, ok := session.Values["user_id"].(int64)
	if !ok {
		return 0, fmt.Errorf("invalid session")
	}
	return userID, nil
}

func destroySession(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "user-session")
	if err != nil {
		return err
	}
	session.Values = map[string]interface{}{}
	return session.Save(r, w)
}
