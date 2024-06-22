package repository

import (
	"TalkBoard/models"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CommentPostgres struct {
	db *sqlx.DB
}

func NewCommentPostgres(db *sqlx.DB) *CommentPostgres {
	return &CommentPostgres{db: db}
}

func (r *CommentPostgres) Create(comment models.Comment) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE id = $1", commentsTable)
	err := r.db.Get(&id, query, comment.ParentCommentId)
	if err != nil {
		var commentId int
		createCommentQuery := fmt.Sprintf("INSERT INTO %s (post_id, user_id, content)"+
			"VALUES ($1, $2, $3) RETURNING id", commentsTable)
		row := r.db.QueryRow(createCommentQuery, comment.PostId, comment.UserId, comment.Content)
		if err := row.Scan(&commentId); err != nil {
			return 0, err
		}

		return commentId, nil
	}

	queryCheck := fmt.Sprintf("SELECT id FROM %s com "+
		" INNER JOIN %s post on com.post_id = post.id "+
		" WHERE com.id = $1 AND post.id = $2", commentsTable, postsTable)
	err = r.db.Get(&id, queryCheck, comment.ParentCommentId, comment.PostId)
	if err != nil {
		return 0, errors.New("postId and commentId do not match each other!")
	}

	createCommentQuery := fmt.Sprintf("INSERT INTO %s (parent_comment_id, post_id, user_id, content)"+
		"VALUES ($1, $2, $3, $4) RETURNING id", commentsTable)
	row := r.db.QueryRow(createCommentQuery, comment.ParentCommentId, comment.PostId, comment.UserId, comment.Content)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
