package repository

import (
	"TalkBoard/models"
	"errors"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"log"
	"testing"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	type mockBehavior func(args models.User, id int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		user         models.User
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			user: models.User{
				Name:     "name",
				Password: "password",
				Email:    "test@mail.ru",
			},
			id: 2,
			mockBehavior: func(args models.User, id int) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(args.Name, args.Email, args.Password).
					WillReturnRows(rows)
			},
		},
		{
			name: "Empty Fields",
			user: models.User{
				Name:  "name",
				Email: "test@mail.ru",
			},
			mockBehavior: func(args models.User, id int) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(0, errors.New("some error"))
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(args.Name, args.Email, args.Password).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.user, testCase.id)

			got, err := r.CreateUser(testCase.user)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}

func TestAuthPostgres_GetUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	type mockBehavior func(email, password string)
	type inputBody struct {
		email    string
		password string
	}

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		user         models.User
		inputBody    inputBody
		wantErr      bool
	}{
		{
			name: "OK",
			user: models.User{
				Id:    1,
				Name:  "name",
				Email: "test@mail.ru",
			},
			inputBody: inputBody{
				email:    "name",
				password: "password",
			},
			mockBehavior: func(email, password string) {
				rows := sqlmock.NewRows([]string{"id", "name", "email"}).
					AddRow(1, "name", "test@mail.ru")

				mock.ExpectQuery("SELECT id, name, email FROM users WHERE (.+) AND (.+)").
					WithArgs(email, password).WillReturnRows(rows)
			},
		},
		{
			name: "Empty Fields",
			user: models.User{
				Id:    1,
				Name:  "name",
				Email: "test@mail.ru",
			},
			inputBody: inputBody{
				email: "name",
			},
			mockBehavior: func(email, password string) {
				rows := sqlmock.NewRows([]string{"id", "name", "email"}).
					AddRow(1, "name", "test@mail.ru").RowError(0, errors.New("some error"))
				mock.ExpectQuery("SELECT id, name, email FROM users WHERE (.+) AND (.+)").
					WithArgs(email, password).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.inputBody.email, testCase.inputBody.password)

			got, err := r.GetUser(testCase.inputBody.email, testCase.inputBody.password)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.user, got)
			}
		})
	}
}
