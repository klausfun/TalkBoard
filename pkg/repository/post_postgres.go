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
	createPostQuery := fmt.Sprintf("INSERT INTO %s (title, content, user_id, access_to_comments)"+
		"VALUES ($1, $2, $3, $4) RETURNING id", postsTable)
	row := r.db.QueryRow(createPostQuery, post.Title, post.Content, userId, post.AccessToComments)
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

func (r *PostPostgres) GetByPostId(postId int) (models.Post, error) {
	var post models.Post
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", postsTable)
	err := r.db.Get(&post, query, postId)

	return post, err
}
