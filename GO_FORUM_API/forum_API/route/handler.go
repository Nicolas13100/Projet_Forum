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
	r.HandleFunc("/getHome", API.GetAllTopic).Methods("GET")
	r.HandleFunc("/getPost/{id}", API.GetTopic).Methods("GET")
	r.HandleFunc("/getUser/{id}", API.GetUser).Methods("GET")
	r.HandleFunc("/getAllTopicMessage/{id}", API.GetAllTopicMessage).Methods("GET")
	r.HandleFunc("/getAllMessageAnswer/{id}", API.GetAllMessageAnswer).Methods("GET")

	// Admin
	r.HandleFunc("/modifyTopic", API.ModifyTopicHandler).Methods("POST")
	r.HandleFunc("/banUser", API.BanUserHandler).Methods("POST")

	// Logging
	r.HandleFunc("/register", API.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", API.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", API.LogoutHandler).Methods("POST")
	r.HandleFunc("/dashboard", API.DashboardHandler).Methods("GET")
	r.HandleFunc("/changeUserData", API.ChangeUserDataHandler).Methods("POST")

	// Create Topic
	r.HandleFunc("/createTopic", API.CreateTopicHandler).Methods("POST")
	r.HandleFunc("/deleteTopic", API.DeleteTopicHandler).Methods("DELETE")
	r.HandleFunc("/deleteComment", API.DeleteCommentHandler).Methods("DELETE")

	// Like Topic
	r.HandleFunc("/likeTopic", API.LikeTopicHandler).Methods("POST")
	r.HandleFunc("/dislikeTopic", API.DislikeTopicHandler).Methods("POST")

	// Favorite Topic
	r.HandleFunc("/favoriteTopic", API.FavoriteTopicHandler).Methods("POST")

	// Comment Topic
	r.HandleFunc("/commentTopic", API.CommentHandler).Methods("POST") // work for both comment a topic and comment a comment
	r.HandleFunc("/updateComment", API.UpdateCommentHandler).Methods("POST")

	// Like Comment
	r.HandleFunc("/likeComment", API.LikeCommentHandler).Methods("POST")
	r.HandleFunc("/dislikeComment", API.DislikeCommentHandler).Methods("POST")

	// UserFriends
	r.HandleFunc("/addFriend", API.AddFriendHandler).Methods("POST")
	r.HandleFunc("/acceptFriend", API.AcceptFriendHandler).Methods("POST")
	r.HandleFunc("/declineFriend", API.DeclineFriendHandler).Methods("POST")
	r.HandleFunc("/deleteFriend", API.DeleteFriendHandler).Methods("POST")

	// Search
	r.HandleFunc("/search", API.SearchHandler).Methods("GET")

	// Serve static files from the "forum_API/static" directory
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("forum_API/static"))))

	// Print statement indicating server is running
	fmt.Println("Server is running on :8080 http://localhost:8080")

	// Start the server
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
