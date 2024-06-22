package repository

import (
	"TalkBoard/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PostPostgres struct {
	db *sqlx.DB
}

func NewPostPostgres(db *sqlx.DB) *PostPostgres {
	return &PostPostgres{db: db}
}

func (r *PostPostgres) Create(userId int, post models.Post) (int, error) {
	var id int
	createEventQuery := fmt.Sprintf("INSERT INTO %s (title, content, user_id, access_to_comments)"+
		"VALUES ($1, $2, $3, $4) RETURNING id", postsTable)
	row := r.db.QueryRow(createEventQuery, post.Title, post.Content, userId, post.AccessToComments)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostPostgres) GetAll() ([]models.Post, error) {
	var posts []models.Post
	query := fmt.Sprintf("SELECT * FROM %s", postsTable)
	err := r.db.Select(&posts, query)

	return posts, err
}
