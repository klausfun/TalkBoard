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

func TestHandler_createComment(t *testing.T) {
	type mockBehavior func(s *mock_service.MockComment, comment models.Comment)

	testTable := []struct {
		name                 string
		inputComment         models.Comment
		inputArgs            map[string]interface{}
		mockBehavior         mockBehavior
		expectedErrorMessage string
		expectedResponseBody models.Comment
	}{
		{
			name: "OK",
			inputArgs: map[string]interface{}{
				"parent_comment_id": 0,
				"userId":            1,
				"postId":            2,
				"content":           "content",
			},
			inputComment: models.Comment{
				ParentCommentId: 0,
				UserId:          1,
				PostId:          2,
				Content:         "content",
			},
			mockBehavior: func(s *mock_service.MockComment, comment models.Comment) {
				s.EXPECT().Create(comment).Return(3, nil)
			},
			expectedErrorMessage: "",
			expectedResponseBody: models.Comment{
				ParentCommentId: 0,
				UserId:          1,
				PostId:          2,
				Content:         "content",
				Id:              3,
			},
		},
		{
			name: "Empty Fields",
			inputArgs: map[string]interface{}{
				"parent_comment_id": 0,
				"userId":            1,
				"postId":            2,
			},
			inputComment: models.Comment{
				ParentCommentId: 0,
				UserId:          1,
				PostId:          2,
			},
			mockBehavior:         func(s *mock_service.MockComment, comment models.Comment) {},
			expectedErrorMessage: "invalid input body",
			expectedResponseBody: models.Comment{},
		},
		{
			name:         "Invalid Input Type",
			inputComment: models.Comment{},
			inputArgs: map[string]interface{}{
				"input": "invalid",
			},
			mockBehavior:         func(s *mock_service.MockComment, comment models.Comment) {},
			expectedErrorMessage: "invalid input body",
			expectedResponseBody: models.Comment{},
		},
		{
			name: "Service Failure",
			inputArgs: map[string]interface{}{
				"parent_comment_id": 0,
				"userId":            1,
				"postId":            2,
				"content":           "content",
			},
			inputComment: models.Comment{
				ParentCommentId: 0,
				UserId:          1,
				PostId:          2,
				Content:         "content",
			},
			mockBehavior: func(s *mock_service.MockComment, comment models.Comment) {
				s.EXPECT().Create(comment).Return(0, errors.New("service failure"))
			},
			expectedErrorMessage: "service failure",
			expectedResponseBody: models.Comment{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCommentService := mock_service.NewMockComment(ctrl)
			testCase.mockBehavior(mockCommentService, testCase.inputComment)

			services := &service.Service{Comment: mockCommentService}
			handler := NewHandler(services)

			args := map[string]interface{}{
				"input": testCase.inputArgs,
			}
			p := graphql.ResolveParams{
				Context: context.Background(),
				Args:    args,
			}

			result, err := handler.createComment(p)
			if err != nil {
				formattedErr, ok := err.(gqlerrors.FormattedError)
				if !ok {
					t.Fatalf("expected gqlerrors.FormattedError, got %T", err)
				}
				assert.Equal(t, testCase.expectedErrorMessage, formattedErr.Message)
				return
			}

			response, ok := result.(models.Comment)
			if !ok {
				t.Fatalf("signUp returned unexpected type: %T", result)
			}

			assert.Equal(t, testCase.expectedResponseBody, response)
		})
	}
}
