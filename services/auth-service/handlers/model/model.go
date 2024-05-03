package model

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ReturnedUser struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type LoggedUserInfo struct {
	RefreshToken string `json:"refresh_token"`
}
