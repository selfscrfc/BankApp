package models

type SignIn struct {
	Login    string `json:"login"`
	Password string `json:"password" validatee:"uuid"`
}

type SignUp struct {
	SignIn
	FullName string `json:"fullname"`
}
