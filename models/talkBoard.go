package models

type Post struct {
	Id               int    `json:"id"`
	UserId           int    `json:"user_id"`
	Title            string `json:"title" binding:"required"`
	Content          string `json:"content" binding:"required"`
	AccessToComments bool   `json:"access_to_comments" binding:"required"`
}

type Subscription struct {
	Id          int    `json:"id"`
	ParentSubId int    `json:"parent_sub_id"`
	PostId      int    `json:"post_id"`
	Content     string `json:"content"`
}
