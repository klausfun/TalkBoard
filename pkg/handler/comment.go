package handler

import (
	"TalkBoard/models"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *Handler) createComment(p graphql.ResolveParams) (interface{}, error) {
	input := p.Args["input"].(map[string]interface{})
	userId := input["userId"].(int)
	postId := input["postId"].(int)
	content := input["content"].(string)

	var parentCommentId int
	if id, ok := input["parentCommentId"]; ok {
		parentCommentId = id.(int)
	} else {
		parentCommentId = 0
	}

	comment := models.Comment{
		UserId:          userId,
		Content:         content,
		PostId:          postId,
		ParentCommentId: parentCommentId,
	}

	commentId, err := h.services.Comment.Create(comment)
	if err != nil {
		return nil, err
	}

	comment.Id = commentId
	return comment, nil
}

func (h *Handler) getSubscriptionsByPostId(c *gin.Context) {}

func (h *Handler) deleteSubscription(c *gin.Context) {}
