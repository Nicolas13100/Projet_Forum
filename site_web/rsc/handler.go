package API

import (
	"fmt"
	"log"
	"net/http"
)

// RUN function sets up HTTP routes and starts the server
func RUN() {
	// Setting up HTTP routes for different endpoints
	http.HandleFunc("/", ErrorHandler)
	http.HandleFunc("/home", indexHandler)

	// Loggin
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/confirmRegister", confirmRegisterHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/successLogin", successLoginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/gestion", gestionHandler)
	http.HandleFunc("/changeLogin", changeLoginHandler)
	//

	// Serve static files from the "site_web/static" directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("site_web/static"))))

	// Print statement indicating server is running
	fmt.Println("Server is running on :8080 http://localhost:8080/home")

	// Start the server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// ErrorHandler handles 404 errors
func ErrorHandler(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "404", logged)
}

// indexHandler handles requests for the home page
func indexHandler(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "home", nil)
}
