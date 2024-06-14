package API

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Mail     string `json:"mail"`
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
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Token    string `json:"token,omitempty"`
	JsonResp []byte `json:"json_resp,omitempty"`
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
