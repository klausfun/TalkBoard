package handler

import (
	"TalkBoard/models"
	"github.com/graphql-go/graphql"
)

func (h *Handler) createPost(p graphql.ResolveParams) (interface{}, error) {
	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return nil, newErrorResponse("invalid input body")
	}

	userId, userIdOk := input["userId"].(int)
	title, titleOk := input["title"].(string)
	content, contentOk := input["content"].(string)
	accessToComments, accessToCommentsOk := input["accessToComments"].(bool)
	if !userIdOk || !titleOk || !accessToCommentsOk || !contentOk || content == "" || title == "" {
		return nil, newErrorResponse("invalid input body")
	}

	post := models.Post{
		Id:               0,
		UserId:           userId,
		Title:            title,
		Content:          content,
		AccessToComments: accessToComments,
	}

	postId, err := h.services.Post.Create(userId, post)
	if err != nil {
		return nil, newErrorResponse("service failure")
	}

	post.Id = postId
	return post, nil
}

func (h *Handler) getAllPosts(p graphql.ResolveParams) (interface{}, error) {
	posts, err := h.services.Post.GetAll()
	if err != nil {
		return nil, newErrorResponse("service failure")
	}

	return posts, nil
}

func (h *Handler) getPostById(p graphql.ResolveParams) (interface{}, error) {
	postId := p.Args["postId"].(int)
	limit, limitOk := p.Args["limit"].(int)
	offset, offsetOk := p.Args["offset"].(int)
	if !limitOk {
		limit = 10
	}
	if !offsetOk {
		offset = 0
	}

	post, err := h.services.Post.GetByPostId(postId)
	if err != nil {
		return nil, err
	}

	comments, err := h.services.Comment.GetByPostId(postId, limit, offset)
	if err != nil {
		return nil, err
	}

	postWithComments := struct {
		Post     models.Post      `json:"post"`
		Comments []models.Comment `json:"comments"`
	}{
		Post:     post,
		Comments: comments,
	}

	return postWithComments, nil
}
