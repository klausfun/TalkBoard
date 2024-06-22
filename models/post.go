package models

type Post struct {
	Id               int    `json:"id" db:"id"`
	UserId           int    `json:"user_id" db:"user_id" binding:"required"`
	Title            string `json:"title" db:"title" binding:"required"`
	Content          string `json:"content" db:"content" binding:"required"`
	AccessToComments bool   `json:"access_to_comments" db:"access_to_comments" binding:"required"`
}
