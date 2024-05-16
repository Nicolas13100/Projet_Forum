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

	//Admin
	http.HandleFunc("/modifyTopic", authenticate(modifyTopicHandler))
	http.HandleFunc("/banUser", authenticate(banUserHandler))

	// Loggin
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", authenticate(logoutHandler))
	http.HandleFunc("/dashboard", authenticate(dashboardHandler))
	http.HandleFunc("/changeUserData", authenticate(changeUserDataHandler))
	//

	// Create Topic
	http.HandleFunc("/createTopic", authenticate(createTopicHandler))
	http.HandleFunc("/deleteTopic", authenticate(deleteTopicHandler))
	http.HandleFunc("/deleteTopicComment", authenticate(deleteTopicCommentHandler))

	// Like Topic
	http.HandleFunc("/likeTopic", authenticate(likeTopicHandler))

	// Favorite Topic
	http.HandleFunc("/favoriteTopic", authenticate(favoriteTopicHandler))

	// Comment Topic
	http.HandleFunc("/commentTopic", authenticate(commentTopicHandler))
	http.HandleFunc("/deleteComment", authenticate(deleteCommentHandler))
	http.HandleFunc("/updateComment", authenticate(updateCommentHandler))

	// Like Comment
	http.HandleFunc("/likeComment", authenticate(likeCommentHandler))
	http.HandleFunc("/dislikeComment", authenticate(dislikeCommentHandler))

	// UserFriends
	http.HandleFunc("/addFriend", authenticate(addFriendHandler))
	http.HandleFunc("/acceptFriend", authenticate(acceptFriendHandler))
	http.HandleFunc("/declineFriend", authenticate(declineFriendHandler))
	http.HandleFunc("/deleteFriend", authenticate(deleteFriendHandler))

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
