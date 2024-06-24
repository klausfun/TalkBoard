package repository

import (
	"TalkBoard/models"
	"errors"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"log"
	"testing"
)

func TestPostPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewPostPostgres(db)

	type args struct {
		userId int
		post   models.Post
	}
	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				userId: 1,
				post: models.Post{
					UserId:           1,
					Title:            "test title",
					Content:          "test content",
					AccessToComments: true,
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO posts").
					WithArgs(args.post.Title, args.post.Content, args.post.UserId, args.post.AccessToComments).
					WillReturnRows(rows)
			},
		},
		{
			name: "Empty Fields",
			args: args{
				userId: 1,
				post: models.Post{
					UserId:           1,
					Title:            "",
					Content:          "test content",
					AccessToComments: true,
				},
			},
			mockBehavior: func(args args, id int) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(0, errors.New("some error"))
				mock.ExpectQuery("INSERT INTO posts").
					WithArgs(args.post.Title, args.post.Content, args.post.UserId, args.post.AccessToComments).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.Create(testCase.args.userId, testCase.args.post)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}
