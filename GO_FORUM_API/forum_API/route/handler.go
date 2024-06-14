package route

import (
	"KoKo/forum_API/rsc"
	"fmt"
	"log"
	"net/http"
)

// RUN function sets up HTTP routes and starts the server
func RUN() {
	// Setting up HTTP routes for different endpoints
	http.HandleFunc("/getHome", API.Authenticate(API.GetAllTopic))
	http.HandleFunc("/getPost/:id", API.Authenticate(API.GetTopic))
	http.HandleFunc("/getUser/:id", API.Authenticate(API.GetUser))

	//Admin
	http.HandleFunc("/modifyTopic", API.Authenticate(API.ModifyTopicHandler))
	http.HandleFunc("/banUser", API.Authenticate(API.BanUserHandler))

	// Loggin
	http.HandleFunc("/register", API.RegisterHandler)
	http.HandleFunc("/login", API.LoginHandler)
	http.HandleFunc("/logout", API.Authenticate(API.LogoutHandler))
	http.HandleFunc("/dashboard", API.Authenticate(API.DashboardHandler))
	http.HandleFunc("/changeUserData", API.Authenticate(API.ChangeUserDataHandler))
	//

	// Create Topic
	http.HandleFunc("/createTopic", API.Authenticate(API.CreateTopicHandler))
	http.HandleFunc("/deleteTopic", API.Authenticate(API.DeleteTopicHandler))
	http.HandleFunc("/deleteComment", API.Authenticate(API.DeleteCommentHandler))

	// Like Topic
	http.HandleFunc("/likeTopic", API.Authenticate(API.LikeTopicHandler))
	http.HandleFunc("/dislikeTopic", API.Authenticate(API.DislikeTopicHandler))

	// Favorite Topic
	http.HandleFunc("/favoriteTopic", API.Authenticate(API.FavoriteTopicHandler))

	// Comment Topic
	http.HandleFunc("/commentTopic", API.Authenticate(API.CommentHandler))
	http.HandleFunc("/updateComment", API.Authenticate(API.UpdateCommentHandler))

	// Like Comment
	http.HandleFunc("/likeComment", API.Authenticate(API.LikeCommentHandler))
	http.HandleFunc("/dislikeComment", API.Authenticate(API.DislikeCommentHandler))

	// UserFriends
	http.HandleFunc("/addFriend", API.Authenticate(API.AddFriendHandler))
	http.HandleFunc("/acceptFriend", API.Authenticate(API.AcceptFriendHandler))
	http.HandleFunc("/declineFriend", API.Authenticate(API.DeclineFriendHandler))
	http.HandleFunc("/deleteFriend", API.Authenticate(API.DeleteFriendHandler))

	// Serve static files from the "forum_API/static" directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("forum_API/static"))))

	// Print statement indicating server is running
	fmt.Println("Server is running on :8080 http://localhost:8080/home")

	// Start the server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
