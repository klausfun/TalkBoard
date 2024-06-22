package models

type Comment struct {
	Id              int       `json:"id" db:"id"`
	UserId          int       `json:"user_id" db:"user_id" binding:"required"`
	ParentCommentId int       `json:"parent_comment_id" db:"parent_comment_id" binding:"required"`
	PostId          int       `json:"post_id" db:"post_id" binding:"required"`
	Content         string    `json:"content" db:"content" binding:"required"`
	Replies         []Comment `json:"replies"`
}
