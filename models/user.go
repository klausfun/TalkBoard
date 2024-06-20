package models

type User struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
