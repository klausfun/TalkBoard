package repository

import (
	"TalkBoard/models"
	"errors"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"log"
	"testing"
)

func TestCommentPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewCommentPostgres(db)

	type mockBehavior func(comment models.Comment, id int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		comment      models.Comment
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			comment: models.Comment{
				UserId:          1,
				ParentCommentId: 2,
				PostId:          3,
				Content:         "test content",
			},
			id: 4,
			mockBehavior: func(comment models.Comment, id int) {
				accessToComments := sqlmock.NewRows([]string{"access_to_comments"}).AddRow(true)
				mock.ExpectQuery("SELECT access_to_comments FROM posts WHERE (.+)").
					WithArgs(comment.PostId).WillReturnRows(accessToComments)

				// случай когда есть родиетль
				commId := sqlmock.NewRows([]string{"id"}).AddRow(2)
				mock.ExpectQuery("SELECT id FROM comments WHERE (.+)").
					WithArgs(comment.ParentCommentId).WillReturnRows(commId)

				parCommId := sqlmock.NewRows([]string{"id"}).AddRow(2)
				mock.ExpectQuery("SELECT com.id FROM comments com INNER JOIN posts post on (.+) "+
					" WHERE (.+) AND (.+)").
					WithArgs(comment.ParentCommentId, comment.PostId).WillReturnRows(parCommId)

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO comments").
					WithArgs(comment.ParentCommentId, comment.PostId, comment.UserId, comment.Content).
					WillReturnRows(rows)
			},
		},
		{
			name: "OK",
			comment: models.Comment{
				UserId:          1,
				ParentCommentId: 2,
				PostId:          3,
				Content:         "test content",
			},
			id: 4,
			mockBehavior: func(comment models.Comment, id int) {
				accessToComments := sqlmock.NewRows([]string{"access_to_comments"}).AddRow(true)
				mock.ExpectQuery("SELECT access_to_comments FROM posts WHERE (.+)").
					WithArgs(comment.PostId).WillReturnRows(accessToComments)

				// случай когда нет родиетля
				commId := sqlmock.NewRows([]string{"id"}).AddRow(4).
					RowError(0, errors.New("no ParentCommentId"))
				mock.ExpectQuery("SELECT id FROM comments WHERE (.+)").
					WithArgs(comment.ParentCommentId).WillReturnRows(commId)

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO comments").
					WithArgs(comment.PostId, comment.UserId, comment.Content).
					WillReturnRows(rows)
			},
		},
		{
			name: "Empty Fields",
			comment: models.Comment{
				UserId:          1,
				ParentCommentId: 2,
				Content:         "test content",
			},
			mockBehavior: func(comment models.Comment, id int) {
				accessToComments := sqlmock.NewRows([]string{"access_to_comments"}).AddRow(true).
					RowError(0, errors.New("some error"))
				mock.ExpectQuery("SELECT access_to_comments FROM posts WHERE (.+)").
					WithArgs(comment.PostId).WillReturnRows(accessToComments)
			},
			wantErr: true,
		},
		{
			name: "There is no such post",
			comment: models.Comment{
				UserId:          1,
				ParentCommentId: 2,
				PostId:          3,
				Content:         "test content",
			},
			mockBehavior: func(comment models.Comment, id int) {
				accessToComments := sqlmock.NewRows([]string{"access_to_comments"}).AddRow(true).
					RowError(0, errors.New("there is no post with this id!"))
				mock.ExpectQuery("SELECT access_to_comments FROM posts WHERE (.+)").
					WithArgs(comment.PostId).WillReturnRows(accessToComments)
			},
			wantErr: true,
		},
		{
			name: "There is no access",
			comment: models.Comment{
				UserId:          1,
				ParentCommentId: 2,
				PostId:          3,
				Content:         "test content",
			},
			mockBehavior: func(comment models.Comment, id int) {
				accessToComments := sqlmock.NewRows([]string{"access_to_comments"}).AddRow(false)
				mock.ExpectQuery("SELECT access_to_comments FROM posts WHERE (.+)").
					WithArgs(comment.PostId).WillReturnRows(accessToComments)
			},
			wantErr: true,
		},
		{
			name: "Incorrect Fields",
			comment: models.Comment{
				UserId:          1,
				ParentCommentId: 2,
				PostId:          3,
				Content:         "test content",
			},
			mockBehavior: func(comment models.Comment, id int) {
				accessToComments := sqlmock.NewRows([]string{"access_to_comments"}).AddRow(true)
				mock.ExpectQuery("SELECT access_to_comments FROM posts WHERE (.+)").
					WithArgs(comment.PostId).WillReturnRows(accessToComments)

				// случай когда есть родиетль
				commId := sqlmock.NewRows([]string{"id"}).AddRow(2)
				mock.ExpectQuery("SELECT id FROM comments WHERE (.+)").
					WithArgs(comment.ParentCommentId).WillReturnRows(commId)

				parCommId := sqlmock.NewRows([]string{"id"}).AddRow(2).
					RowError(0, errors.New("postId and commentId do not match each other"))
				mock.ExpectQuery("SELECT com.id FROM comments com INNER JOIN posts post on (.+) "+
					" WHERE (.+) AND (.+)").
					WithArgs(comment.ParentCommentId, comment.PostId).WillReturnRows(parCommId)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.comment, testCase.id)

			got, err := r.Create(testCase.comment)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}

func TestCommentPostgres_GetByPostId(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewCommentPostgres(db)

	type inputBody struct {
		postId, limit, offset int
	}
	type mockBehavior func(args inputBody)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         inputBody
		comments     []models.Comment
		wantErr      bool
	}{
		{
			name: "OK",
			comments: []models.Comment{
				{
					Id:      1,
					UserId:  1,
					PostId:  3,
					Content: "content1",
				},
			},
			args: inputBody{
				postId: 3,
				limit:  2,
				offset: 0,
			},
			mockBehavior: func(args inputBody) {
				comments := sqlmock.NewRows([]string{"id", "post_id", "user_id", "content"}).
					AddRow(1, 3, 1, "content1")
				mock.ExpectQuery("SELECT id, post_id, user_id, content "+
					" FROM comments WHERE post_id = (.+) AND parent_comment_id IS NULL"+
					" ORDER BY id LIMIT (.+) OFFSET (.+)").
					WithArgs(args.postId, args.limit, args.offset).WillReturnRows(comments)

				replies := sqlmock.NewRows([]string{"id", "post_id", "user_id", "parent_comment_id", "content"})
				mock.ExpectQuery("SELECT id, post_id, user_id, parent_comment_id, content " +
					"FROM comments WHERE parent_comment_id = (.+)").
					WithArgs(1).WillReturnRows(replies)
			},
		},
		{
			name: "There is no such comment",
			args: inputBody{
				postId: 3,
				limit:  2,
				offset: 0,
			},
			mockBehavior: func(args inputBody) {
				comments := sqlmock.NewRows([]string{"id", "post_id", "user_id", "content"}).
					AddRow(1, 3, 1, "content1").RowError(0, errors.New("error"))
				mock.ExpectQuery("SELECT id, post_id, user_id, content "+
					" FROM comments WHERE post_id = (.+) AND parent_comment_id IS NULL"+
					" ORDER BY id LIMIT (.+) OFFSET (.+)").
					WithArgs(args.postId, args.limit, args.offset).WillReturnRows(comments)
			},
			wantErr: true,
		},
		{
			name: "There is no such comment",
			args: inputBody{
				postId: 3,
				limit:  2,
				offset: 0,
			},
			mockBehavior: func(args inputBody) {
				comments := sqlmock.NewRows([]string{"id", "post_id", "user_id", "content"}).
					AddRow(1, 3, 1, "content1")
				mock.ExpectQuery("SELECT id, post_id, user_id, content "+
					" FROM comments WHERE post_id = (.+) AND parent_comment_id IS NULL"+
					" ORDER BY id LIMIT (.+) OFFSET (.+)").
					WithArgs(args.postId, args.limit, args.offset).WillReturnRows(comments)

				replies := sqlmock.NewRows([]string{"id", "post_id", "user_id", "parent_comment_id", "content"}).
					RowError(0, errors.New("error"))
				mock.ExpectQuery("SELECT id, post_id, user_id, parent_comment_id, content " +
					"FROM comments WHERE parent_comment_id = (.+)").
					WithArgs(0).WillReturnRows(replies)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.GetByPostId(testCase.args.postId, testCase.args.limit, testCase.args.offset)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.comments, got)
			}
		})
	}
}
