package handler

import (
	"TalkBoard/models"
	"github.com/graphql-go/graphql"
)

func (h *Handler) signUp(p graphql.ResolveParams) (interface{}, error) {
	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return nil, newErrorResponse("invalid input body")
	}
	name, nameOk := input["name"].(string)
	email, emailOk := input["email"].(string)
	password, passwordOk := input["password"].(string)

	if !nameOk || !emailOk || !passwordOk || name == "" || email == "" || password == "" {
		return nil, newErrorResponse("invalid input body")
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	id, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		return nil, newErrorResponse("service failure")
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
