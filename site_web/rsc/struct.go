package API

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Mail     string `json:"mail"`
}

// UserData structure for individual user data
type UserData struct {
	Fav []int `json:"fav"`
}

// struct for search
type GameInfo struct {
	GameId int `json:"game"`
}

type CombinedData struct {
	Result interface{}
	Name   string
	Logged bool
}
