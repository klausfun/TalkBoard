package handler

import (
	"TalkBoard/models"
	"github.com/graphql-go/graphql"
)

func (h *Handler) createComment(p graphql.ResolveParams) (interface{}, error) {
	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return nil, newErrorResponse("invalid input body")
	}

	userId, userIdOk := input["userId"].(int)
	postId, postIdOk := input["postId"].(int)
	content, contentOk := input["content"].(string)
	if !userIdOk || !postIdOk || !contentOk || content == "" {
		return nil, newErrorResponse("invalid input body")
	}

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
		return nil, newErrorResponse("service failure")
	}

	comment.Id = commentId
	return comment, nil
}

func (h *Handler) deleteComment(p graphql.ResolveParams) {}
