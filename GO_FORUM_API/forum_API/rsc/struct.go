package API

import "time"

// User struct to represent user data
type User struct {
	UserID           int    `json:"user_id"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	RegistrationDate string `json:"registration_date"`
	LastLoginDate    string `json:"last_login_date"`
	Biography        string `json:"biography"`
	IsAdmin          bool   `json:"is_admin"`
	IsModerator      bool   `json:"is_moderator"`
	IsDeleted        bool   `json:"is_deleted"`
	ProfilePic       string `json:"profile_pic"`
}

// UserData structure for individual user data
type UserData struct {
	Fav []int `json:"fav"`
}

// GameInfo struct for search
type GameInfo struct {
	GameId int `json:"game"`
}

type CombinedData struct {
	Result interface{}
	Name   string
	Logged bool
}

// Response struct for API response
type Response struct {
	Message string `json:"message"`
}

// APIResponse represents the structure of the API response
type APIResponse struct {
	Status       int                    `json:"status"`
	Message      string                 `json:"message"`
	Token        string                 `json:"token,omitempty"`
	Resp         map[string]interface{} `json:"resp,omitempty"`
	JsonResp     []byte                 `json:"json_resp,omitempty"`
	Data         []TagInfo              `json:"data,omitempty"`
	UserData     map[string]string      `json:"UserData,omitempty"`
	NumberOfLike map[string]int         `json:"NumberOfLike,omitempty"`
	UsersData    []User                 `json:"UsersData,omitempty"`
	FollowerData map[string]int         `json:"FollowerData,omitempty"`
	UserID       int                    `json:"UserID,omitempty"`
	User         User                   `json:"user,omitempty"`
}

type Topic struct {
	TopicID      int    `json:"topic_id"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	CreationDate string `json:"creation_date"`
	Status       int    `json:"status"`
	IsPrivate    bool   `json:"is_private"`
	UserID       int    `json:"user_id"`
}

// Message struct to represent a message in the database
type Message struct {
	MessageID     int       `json:"message_id"`
	Body          string    `json:"body"`
	DateSent      time.Time `json:"date_sent"`
	TopicID       int       `json:"topic_id"`
	BaseMessageID *int      `json:"base_message_id"`
	UserID        int       `json:"user_id"`
}

type SearchResults struct {
	Topics   []Topic   `json:"topics"`
	Messages []Message `json:"messages"`
}

type TagInfo struct {
	TagID   int    `json:"tag_id"`
	TagName string `json:"tag_name"`
}
