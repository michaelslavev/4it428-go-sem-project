package model

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoggedUserInfo struct {
	RefreshToken string `json:"refresh_token"`
}
