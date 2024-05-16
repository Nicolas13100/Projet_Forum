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
