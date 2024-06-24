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

func TestHandler_getAllPosts(t *testing.T) {
	type mockBehavior func(s *mock_service.MockPost)

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedErrorMessage string
		expectedResponseBody []models.Post
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockPost) {
				s.EXPECT().GetAll().Return([]models.Post{
					{
						UserId:           1,
						AccessToComments: true,
						Title:            "title",
						Content:          "content",
						Id:               2,
					},
				}, nil)
			},
			expectedErrorMessage: "",
			expectedResponseBody: []models.Post{
				{
					UserId:           1,
					AccessToComments: true,
					Title:            "title",
					Content:          "content",
					Id:               2,
				},
			},
		},
		{
			name: "Service Failure",
			mockBehavior: func(s *mock_service.MockPost) {
				s.EXPECT().GetAll().Return([]models.Post{}, errors.New("service failure"))
			},
			expectedErrorMessage: "service failure",
			expectedResponseBody: []models.Post{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPostService := mock_service.NewMockPost(ctrl)
			testCase.mockBehavior(mockPostService)

			services := &service.Service{Post: mockPostService}
			handler := NewHandler(services)

			p := graphql.ResolveParams{
				Context: context.Background(),
			}

			result, err := handler.getAllPosts(p)
			if err != nil {
				formattedErr, ok := err.(gqlerrors.FormattedError)
				if !ok {
					t.Fatalf("expected gqlerrors.FormattedError, got %T", err)
				}
				assert.Equal(t, testCase.expectedErrorMessage, formattedErr.Message)
				return
			}

			response, ok := result.([]models.Post)
			if !ok {
				t.Fatalf("signUp returned unexpected type: %T", result)
			}

			assert.Equal(t, testCase.expectedResponseBody, response)
		})
	}
}

func TestHandler_getPostById(t *testing.T) {
	type mockBehaviorPost func(s *mock_service.MockPost, postId int)
	type mockBehaviorComment func(s *mock_service.MockComment, postId, limit, offset int)

	type inputBody struct {
		postId, limit, offset int
	}

	testTable := []struct {
		name                 string
		inputBody            inputBody
		postId               int
		args                 map[string]interface{}
		mockBehaviorPost     mockBehaviorPost
		mockBehaviorComment  mockBehaviorComment
		expectedErrorMessage string
		expectedResponseBody PostWithComments
	}{
		{
			name: "OK",
			inputBody: inputBody{
				postId: 2,
				limit:  10,
				offset: 0,
			},
			postId: 2,
			args: map[string]interface{}{
				"postId": 2,
				"limit":  10,
				"offset": 0,
			},
			mockBehaviorPost: func(s *mock_service.MockPost, postId int) {
				s.EXPECT().GetByPostId(postId).Return(models.Post{
					UserId:           1,
					AccessToComments: true,
					Title:            "title",
					Content:          "content",
					Id:               2,
				}, nil)
			},
			mockBehaviorComment: func(s *mock_service.MockComment, postId, limit, offset int) {
				s.EXPECT().GetByPostId(postId, limit, offset).Return([]models.Comment{
					{
						UserId:          1,
						ParentCommentId: 0,
						PostId:          2,
						Content:         "content",
						Id:              1,
					},
				}, nil)
			},
			expectedErrorMessage: "",
			expectedResponseBody: PostWithComments{
				Post: models.Post{
					UserId:           1,
					AccessToComments: true,
					Title:            "title",
					Content:          "content",
					Id:               2,
				},
				Comments: []models.Comment{
					{
						UserId:          1,
						ParentCommentId: 0,
						PostId:          2,
						Content:         "content",
						Id:              1,
					},
				},
			},
		},
		{
			name: "Empty Fields",
			inputBody: inputBody{
				limit:  10,
				offset: 0,
			},
			args: map[string]interface{}{
				"limit":  10,
				"offset": 0,
			},
			postId:               2,
			mockBehaviorPost:     func(s *mock_service.MockPost, postId int) {},
			mockBehaviorComment:  func(s *mock_service.MockComment, postId, limit, offset int) {},
			expectedErrorMessage: "invalid input body",
			expectedResponseBody: PostWithComments{},
		},
		{
			name: "Service Failure",
			inputBody: inputBody{
				postId: 2,
				limit:  10,
				offset: 0,
			},
			postId: 2,
			args: map[string]interface{}{
				"postId": 2,
				"limit":  10,
				"offset": 0,
			},
			mockBehaviorPost: func(s *mock_service.MockPost, postId int) {
				s.EXPECT().GetByPostId(postId).Return(models.Post{}, errors.New("service failure"))
			},
			mockBehaviorComment:  func(s *mock_service.MockComment, postId, limit, offset int) {},
			expectedErrorMessage: "service failure",
			expectedResponseBody: PostWithComments{},
		},
		{
			name: "Service Failure",
			inputBody: inputBody{
				postId: 2,
				limit:  10,
				offset: 0,
			},
			postId: 2,
			args: map[string]interface{}{
				"postId": 2,
				"limit":  10,
				"offset": 0,
			},
			mockBehaviorPost: func(s *mock_service.MockPost, postId int) {
				s.EXPECT().GetByPostId(postId).Return(models.Post{
					UserId:           1,
					AccessToComments: true,
					Title:            "title",
					Content:          "content",
					Id:               2,
				}, nil)
			},
			mockBehaviorComment: func(s *mock_service.MockComment, postId, limit, offset int) {
				s.EXPECT().GetByPostId(postId, limit, offset).Return(nil, errors.New("service failure"))
			},
			expectedErrorMessage: "service failure",
			expectedResponseBody: PostWithComments{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPostService := mock_service.NewMockPost(ctrl)
			testCase.mockBehaviorPost(mockPostService, testCase.postId)

			mockCommentService := mock_service.NewMockComment(ctrl)
			testCase.mockBehaviorComment(mockCommentService, testCase.inputBody.postId,
				testCase.inputBody.limit, testCase.inputBody.offset)

			services := &service.Service{Post: mockPostService, Comment: mockCommentService}
			handlers := NewHandler(services)

			p := graphql.ResolveParams{
				Context: context.Background(),
				Args:    testCase.args,
			}

			result, err := handlers.getPostById(p)
			if err != nil {
				formattedErr, ok := err.(gqlerrors.FormattedError)
				if !ok {
					t.Fatalf("expected gqlerrors.FormattedError, got %T", err)
				}
				assert.Equal(t, testCase.expectedErrorMessage, formattedErr.Message)
				return
			}

			response, ok := result.(PostWithComments)
			if !ok {
				t.Fatalf("returned unexpected type: %T", result)
			}

			assert.Equal(t, testCase.expectedResponseBody, response)
		})
	}
}
