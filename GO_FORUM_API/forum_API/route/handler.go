package route

import (
	API "KoKo/forum_API/rsc"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RUN() {
	r := mux.NewRouter()

	// Setting up HTTP routes for different endpoints
	r.HandleFunc("/getHome", API.Authenticate(API.GetAllTopic)).Methods("GET")
	r.HandleFunc("/getPost/{id}", API.Authenticate(API.GetTopic)).Methods("GET")
	r.HandleFunc("/getUser/{id}", API.Authenticate(API.GetUser)).Methods("GET")
	r.HandleFunc("/getAllTopicMessage/{id}", API.Authenticate(API.GetAllTopicMessage)).Methods("GET")
	r.HandleFunc("/getAllMessageAnswer/{id}", API.Authenticate(API.GetAllMessageAnswer)).Methods("GET")

	// Admin
	r.HandleFunc("/modifyTopic", API.Authenticate(API.ModifyTopicHandler)).Methods("POST")
	r.HandleFunc("/banUser", API.Authenticate(API.BanUserHandler)).Methods("POST")

	// Logging
	r.HandleFunc("/register", API.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", API.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", API.Authenticate(API.LogoutHandler)).Methods("POST")
	r.HandleFunc("/dashboard", API.Authenticate(API.DashboardHandler)).Methods("GET")
	r.HandleFunc("/changeUserData", API.Authenticate(API.ChangeUserDataHandler)).Methods("POST")

	// Create Topic
	r.HandleFunc("/createTopic", API.Authenticate(API.CreateTopicHandler)).Methods("POST")
	r.HandleFunc("/deleteTopic", API.Authenticate(API.DeleteTopicHandler)).Methods("DELETE")
	r.HandleFunc("/deleteComment", API.Authenticate(API.DeleteCommentHandler)).Methods("DELETE")

	// Like Topic
	r.HandleFunc("/likeTopic", API.Authenticate(API.LikeTopicHandler)).Methods("POST")
	r.HandleFunc("/dislikeTopic", API.Authenticate(API.DislikeTopicHandler)).Methods("POST")

	// Favorite Topic
	r.HandleFunc("/favoriteTopic", API.Authenticate(API.FavoriteTopicHandler)).Methods("POST")

	// Comment Topic
	r.HandleFunc("/commentTopic", API.Authenticate(API.CommentHandler)).Methods("POST") // work for both comment a topic and comment a comment
	r.HandleFunc("/updateComment", API.Authenticate(API.UpdateCommentHandler)).Methods("POST")

	// Like Comment
	r.HandleFunc("/likeComment", API.Authenticate(API.LikeCommentHandler)).Methods("POST")
	r.HandleFunc("/dislikeComment", API.Authenticate(API.DislikeCommentHandler)).Methods("POST")

	// UserFriends
	r.HandleFunc("/addFriend", API.Authenticate(API.AddFriendHandler)).Methods("POST")
	r.HandleFunc("/acceptFriend", API.Authenticate(API.AcceptFriendHandler)).Methods("POST")
	r.HandleFunc("/declineFriend", API.Authenticate(API.DeclineFriendHandler)).Methods("POST")
	r.HandleFunc("/deleteFriend", API.Authenticate(API.DeleteFriendHandler)).Methods("POST")

	// Search
	r.HandleFunc("/search", API.Authenticate(API.SearchHandler)).Methods("GET")

	// Serve static files from the "forum_API/static" directory
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("forum_API/static"))))

	// Print statement indicating server is running
	fmt.Println("Server is running on :8080 http://localhost:8080")

	// Start the server
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
