package model

type User struct {
	UserId   int    `json:"user_id"`
	Username string `json:"user_name"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}
