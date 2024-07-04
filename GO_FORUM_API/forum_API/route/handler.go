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
	r.HandleFunc("/api/getHome/{page}/{pageSize}", API.GetAllTopic).Methods("GET")
	r.HandleFunc("/api/getTopicTag/{id}", API.GetTopicTagsNamesByTopicId).Methods("GET")
	r.HandleFunc("/api/getTopic/{id}", API.GetTopic).Methods("GET")
	r.HandleFunc("/api/getUser/{id}", API.GetUser).Methods("GET")
	r.HandleFunc("/api/getAllTopicMessage/{id}", API.GetAllTopicMessage).Methods("GET")
	r.HandleFunc("/api/getAllMessageAnswer/{id}", API.GetAllMessageAnswer).Methods("GET")
	r.HandleFunc("/api/getTopicOwner/{id}", API.GetTopicOwner).Methods("GET")
	r.HandleFunc("/api/getLikeTopicNumber/{id}", API.GetLikeNumberOfTopic).Methods("GET")
	r.HandleFunc("/api/getRandomUser", API.GetForYouUser).Methods("GET")
	r.HandleFunc("/api/getFollowers/{id}", API.GetUsersFollowers).Methods("GET")
	r.HandleFunc("/api/getUserIDByToken/{token}", API.GetUserIdByToken).Methods("GET")
	r.HandleFunc("/api/logout/{token}", API.DeleteToken).Methods("DELETE")
	r.HandleFunc("/api/getTopicImg/{id}", API.GetTopicImg).Methods("GET")
	r.HandleFunc("/api/getUserTopics/{id}", API.GetUserTopics).Methods("GET")
	r.HandleFunc("/api/getUserFollowings/{id}", API.GetUsersFollowing).Methods("GET")
	r.HandleFunc("/api/getUserFollow/{id}", API.GetUsersFollow).Methods("GET")
	r.HandleFunc("/api/isFollowed/{myId}/{otherId}", API.IsFollower).Methods("GET")
	r.HandleFunc("/api/follow/{userId}/{toFollowUserId}", API.FollowUser).Methods("POST")
	r.HandleFunc("/api/unfollow/{userId}/{toFollowUserId}", API.UnfollowUser).Methods("DELETE")

	// Admin
	r.HandleFunc("/api/modifyTopic", API.ModifyTopicHandler).Methods("POST")
	r.HandleFunc("/api/banUser", API.BanUserHandler).Methods("POST")

	// Logging
	r.HandleFunc("/api/register", API.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", API.LoginHandler).Methods("POST")
	r.HandleFunc("/api/logout", API.LogoutHandler).Methods("POST")
	r.HandleFunc("/api/dashboard", API.DashboardHandler).Methods("GET")
	r.HandleFunc("/api/changeUserData", API.ChangeUserDataHandler).Methods("POST")

	// Create Topic
	r.HandleFunc("/api/createTopic", API.CreateTopicHandler).Methods("POST")
	r.HandleFunc("/api/deleteTopic", API.DeleteTopicHandler).Methods("DELETE")
	r.HandleFunc("/api/deleteComment", API.DeleteCommentHandler).Methods("DELETE")

	// Like Topic
	r.HandleFunc("/api/likeTopic", API.LikeTopicHandler).Methods("POST")
	r.HandleFunc("/api/dislikeTopic", API.DislikeTopicHandler).Methods("POST")

	// Favorite Topic
	r.HandleFunc("/api/favoriteTopic", API.FavoriteTopicHandler).Methods("POST")

	// Comment Topic
	r.HandleFunc("/api/commentTopic", API.CommentHandler).Methods("POST") // work for both comment a topic and comment a comment
	r.HandleFunc("/api/updateComment", API.UpdateCommentHandler).Methods("POST")

	// Like Comment
	r.HandleFunc("/api/likeComment", API.LikeCommentHandler).Methods("POST")
	r.HandleFunc("/api/dislikeComment", API.DislikeCommentHandler).Methods("POST")

	// UserFriends
	r.HandleFunc("/api/addFriend", API.AddFriendHandler).Methods("POST")
	r.HandleFunc("/api/acceptFriend", API.AcceptFriendHandler).Methods("POST")
	r.HandleFunc("/api/declineFriend", API.DeclineFriendHandler).Methods("POST")
	r.HandleFunc("/api/deleteFriend", API.DeleteFriendHandler).Methods("POST")

	// Search
	r.HandleFunc("/api/search", API.SearchHandler).Methods("GET")

	// Serve static files from the "forum_API/static" directory
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("forum_API/static"))))

	// Print statement indicating server is running
	fmt.Println("Server is running on :8080 http://localhost:8080")

	// Start the server
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
