package handler

import (
	"TalkBoard/models"
	"github.com/graphql-go/graphql"
)

func (h *Handler) createPost(p graphql.ResolveParams) (interface{}, error) {
	input := p.Args["input"].(map[string]interface{})
	userId := input["userId"].(int)
	title := input["title"].(string)
	content := input["content"].(string)
	accessToComments := input["accessToComments"].(bool)

	post := models.Post{
		Title:            title,
		Content:          content,
		AccessToComments: accessToComments,
	}

	postId, err := h.services.Post.Create(userId, post)
	if err != nil {
		return nil, err
	}

	post.Id = postId
	return post, nil
}

//func (h *Handler) getAllPosts(c *gin.Context) {}
//
//func (h *Handler) getPostById(c *gin.Context) {}
