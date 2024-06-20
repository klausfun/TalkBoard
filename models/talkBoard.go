package models

type Post struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type Subscription struct {
	Id          int    `json:"id"`
	ParentSubId int    `json:"parent_sub_id"`
	PostId      int    `json:"post_id"`
	Content     string `json:"content"`
}
