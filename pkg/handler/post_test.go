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

func TestHandler_createPost(t *testing.T) {
	type mockBehavior func(s *mock_service.MockPost, post models.Post)

	testTable := []struct {
		name                 string
		inputPost            models.Post
		inputArgs            map[string]interface{}
		mockBehavior         mockBehavior
		expectedErrorMessage string
		expectedResponseBody models.Post
	}{
		{
			name: "OK",
			inputArgs: map[string]interface{}{
				"userId":           1,
				"accessToComments": true,
				"title":            "title",
				"content":          "content",
			},
			inputPost: models.Post{
				UserId:           1,
				AccessToComments: true,
				Title:            "title",
				Content:          "content",
			},
			mockBehavior: func(s *mock_service.MockPost, post models.Post) {
				s.EXPECT().Create(post.UserId, post).Return(2, nil)
			},
			expectedErrorMessage: "",
			expectedResponseBody: models.Post{
				UserId:           1,
				AccessToComments: true,
				Title:            "title",
				Content:          "content",
				Id:               2,
			},
		},
		{
			name: "Empty Fields",
			inputArgs: map[string]interface{}{
				"userId":           1,
				"accessToComments": true,
				"title":            "title",
			},
			inputPost: models.Post{
				UserId:           1,
				AccessToComments: true,
				Title:            "title",
			},
			mockBehavior:         func(s *mock_service.MockPost, post models.Post) {},
			expectedErrorMessage: "invalid input body",
			expectedResponseBody: models.Post{},
		},
		{
			name:      "Invalid Input Type",
			inputPost: models.Post{},
			inputArgs: map[string]interface{}{
				"input": "invalid",
			},
			mockBehavior:         func(s *mock_service.MockPost, post models.Post) {},
			expectedErrorMessage: "invalid input body",
			expectedResponseBody: models.Post{},
		},
		{
			name: "Service Failure",
			inputArgs: map[string]interface{}{
				"userId":           1,
				"accessToComments": true,
				"title":            "title",
				"content":          "content",
			},
			inputPost: models.Post{
				UserId:           1,
				AccessToComments: true,
				Title:            "title",
				Content:          "content",
			},
			mockBehavior: func(s *mock_service.MockPost, post models.Post) {
				s.EXPECT().Create(post.UserId, post).Return(0, errors.New("service failure"))
			},
			expectedErrorMessage: "service failure",
			expectedResponseBody: models.Post{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPostService := mock_service.NewMockPost(ctrl)
			testCase.mockBehavior(mockPostService, testCase.inputPost)

			services := &service.Service{Post: mockPostService}
			handler := NewHandler(services)

			args := map[string]interface{}{
				"input": testCase.inputArgs,
			}
			p := graphql.ResolveParams{
				Context: context.Background(),
				Args:    args,
			}

			result, err := handler.createPost(p)
			if err != nil {
				formattedErr, ok := err.(gqlerrors.FormattedError)
				if !ok {
					t.Fatalf("expected gqlerrors.FormattedError, got %T", err)
				}
				assert.Equal(t, testCase.expectedErrorMessage, formattedErr.Message)
				return
			}

			response, ok := result.(models.Post)
			if !ok {
				t.Fatalf("signUp returned unexpected type: %T", result)
			}

			assert.Equal(t, testCase.expectedResponseBody, response)
		})
	}
}
