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

func TestPostPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewPostPostgres(db)

	testTable := []struct {
		name    string
		mock    func()
		want    []models.Post
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "access_to_comments"}).
					AddRow(1, 1, "title1", "content1", true).
					AddRow(2, 1, "title2", "content2", false).
					AddRow(3, 1, "title3", "content3", true)

				mock.ExpectQuery("SELECT \\* FROM posts").WillReturnRows(rows)
			},
			want: []models.Post{
				{1, 1, "title1", "content1", true},
				{2, 1, "title2", "content2", false},
				{3, 1, "title3", "content3", true},
			},
		},
		{
			name: "No Records",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "access_to_comments"})

				mock.ExpectQuery("SELECT \\* FROM posts").WillReturnRows(rows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, err := r.GetAll()
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, got)
			}
		})
	}
}

func TestPostPostgres_GetByPostId(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewPostPostgres(db)

	testTable := []struct {
		name    string
		mock    func()
		postId  int
		want    models.Post
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "access_to_comments"}).
					AddRow(1, 1, "title", "content", true)

				mock.ExpectQuery("SELECT \\* FROM posts WHERE (.+)").WithArgs(1).WillReturnRows(rows)
			},
			postId: 1,
			want:   models.Post{1, 1, "title", "content", true},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "access_to_comments"})

				mock.ExpectQuery("SELECT \\* FROM posts WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)
			},
			postId:  1,
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, err := r.GetByPostId(testCase.postId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, got)
			}
		})
	}
}
