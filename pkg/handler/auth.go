package handler

import (
	"TalkBoard/models"
	"github.com/graphql-go/graphql"
)

func (h *Handler) signUp(p graphql.ResolveParams) (interface{}, error) {
	input := p.Args["input"].(map[string]interface{})
	user := models.User{
		Name:     input["name"].(string),
		Email:    input["email"].(string),
		Password: input["password"].(string),
	}

	id, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		return nil, err
	}

	user.Id = id
	return user, nil
}

func (h *Handler) signIn(p graphql.ResolveParams) (interface{}, error) {
	input := p.Args["input"].(map[string]interface{})
	email := input["email"].(string)
	password := input["password"].(string)

	token, err := h.services.Authorization.GenerateToken(email, password)
	if err != nil {
		return "", err
	}

	return map[string]interface{}{
		"token": token,
	}, nil
}
