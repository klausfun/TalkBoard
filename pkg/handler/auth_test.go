package handler

import (
	"TalkBoard/models"
	"TalkBoard/pkg/service"
	mock_service "TalkBoard/pkg/service/mocks"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user models.User)

	testTable := []struct {
		name                 string
		inputUser            models.User
		inputArgs            map[string]interface{}
		mockBehavior         mockBehavior
		expectedErrorMessage string
		expectedResponseBody models.User
	}{
		{
			name: "OK",
			inputArgs: map[string]interface{}{
				"email":    "test@mail.ru",
				"name":     "test",
				"password": "qwerty",
			},
			inputUser: models.User{
				Email:    "test@mail.ru",
				Name:     "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedErrorMessage: "",
			expectedResponseBody: models.User{
				Id:       1,
				Email:    "test@mail.ru",
				Name:     "test",
				Password: "qwerty",
			},
		},
		{
			name: "Empty Fields",
			inputUser: models.User{
				Email: "test@mail.ru",
				Name:  "test",
			},
			inputArgs: map[string]interface{}{
				"email": "test@mail.ru",
				"name":  "",
			},
			mockBehavior:         func(s *mock_service.MockAuthorization, user models.User) {},
			expectedErrorMessage: "invalid input body",
			expectedResponseBody: models.User{},
		},
		{
			name:      "Invalid Input Type",
			inputUser: models.User{},
			inputArgs: map[string]interface{}{
				"input": "invalid",
			},
			mockBehavior:         func(s *mock_service.MockAuthorization, user models.User) {},
			expectedErrorMessage: "invalid input body",
			expectedResponseBody: models.User{},
		},
		{
			name: "Service Failure",
			inputArgs: map[string]interface{}{
				"email":    "test@mail.ru",
				"name":     "test",
				"password": "qwerty",
			},
			inputUser: models.User{
				Email:    "test@mail.ru",
				Name:     "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedErrorMessage: "service failure",
			expectedResponseBody: models.User{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mock_service.NewMockAuthorization(ctrl)
			testCase.mockBehavior(mockAuthService, testCase.inputUser)

			services := &service.Service{Authorization: mockAuthService}
			handler := NewHandler(services)

			args := map[string]interface{}{
				"input": testCase.inputArgs,
			}
			p := graphql.ResolveParams{
				Context: context.Background(),
				Args:    args,
			}

			result, err := handler.signUp(p)
			if err != nil {
				formattedErr, ok := err.(gqlerrors.FormattedError)
				if !ok {
					t.Fatalf("expected gqlerrors.FormattedError, got %T", err)
				}
				assert.Equal(t, testCase.expectedErrorMessage, formattedErr.Message)
				return
			}

			response, ok := result.(models.User)
			if !ok {
				t.Fatalf("signUp returned unexpected type: %T", result)
			}

			assert.Equal(t, testCase.expectedResponseBody, response)

		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user models.User)

	testTable := []struct {
		name                 string
		inputUser            models.User
		inputArgs            map[string]interface{}
		mockBehavior         mockBehavior
		expectedErrorMessage string
		expectedResponseBody map[string]interface{}
	}{
		{
			name: "OK",
			inputArgs: map[string]interface{}{
				"email":    "test@mail.ru",
				"password": "qwerty",
			},
			inputUser: models.User{
				Email:    "test@mail.ru",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.User) {
				s.EXPECT().GenerateToken(user.Email, user.Password).Return("1", nil)
			},
			expectedErrorMessage: "",
			expectedResponseBody: map[string]interface{}{
				"token": "1",
			},
		},
		{
			name: "Empty Fields",
			inputUser: models.User{
				Email: "test@mail.ru",
			},
			inputArgs: map[string]interface{}{
				"email": "test@mail.ru",
			},
			mockBehavior:         func(s *mock_service.MockAuthorization, user models.User) {},
			expectedErrorMessage: "invalid input body",
			expectedResponseBody: map[string]interface{}{},
		},
		{
			name:      "Invalid Input Type",
			inputUser: models.User{},
			inputArgs: map[string]interface{}{
				"input": "invalid",
			},
			mockBehavior:         func(s *mock_service.MockAuthorization, user models.User) {},
			expectedErrorMessage: "invalid input body",
			expectedResponseBody: map[string]interface{}{},
		},
		{
			name: "Service Failure",
			inputArgs: map[string]interface{}{
				"email":    "test@mail.ru",
				"password": "qwerty",
			},
			inputUser: models.User{
				Email:    "test@mail.ru",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.User) {
				s.EXPECT().GenerateToken(user.Email, user.Password).Return("", errors.New("service failure"))
			},
			expectedErrorMessage: "service failure",
			expectedResponseBody: map[string]interface{}{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mock_service.NewMockAuthorization(ctrl)
			testCase.mockBehavior(mockAuthService, testCase.inputUser)

			services := &service.Service{Authorization: mockAuthService}
			handler := NewHandler(services)

			args := map[string]interface{}{
				"input": testCase.inputArgs,
			}
			p := graphql.ResolveParams{
				Context: context.Background(),
				Args:    args,
			}

			result, err := handler.signIn(p)
			if err != nil {
				formattedErr, ok := err.(gqlerrors.FormattedError)
				if !ok {
					t.Fatalf("expected gqlerrors.FormattedError, got %T", err)
				}
				assert.Equal(t, testCase.expectedErrorMessage, formattedErr.Message)
				return
			}

			response, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("signIn returned unexpected type: %T", result)
			}

			assert.Equal(t, testCase.expectedResponseBody, response)

		})
	}
}
