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
	var accessToComments bool
	queryCheckAccessToComments := fmt.Sprintf("SELECT access_to_comments FROM %s WHERE id = $1", postsTable)
	err := r.db.Get(&accessToComments, queryCheckAccessToComments, comment.PostId)
	if err != nil {
		return 0, errors.New("there is no post with this id!")
	}
	if !accessToComments {
		return 0, errors.New("no comments are available under this post!")
	}

	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE id = $1", commentsTable)
	err = r.db.Get(&id, query, comment.ParentCommentId)
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

	queryCheck := fmt.Sprintf("SELECT com.id FROM %s com "+
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

func (r *CommentPostgres) GetByPostId(postId, limit, offset int) ([]models.Comment, error) {
	var comments []models.Comment
	query := fmt.Sprintf("SELECT id, post_id, user_id, content FROM %s WHERE post_id = $1 AND parent_comment_id IS NULL"+
		" ORDER BY id LIMIT $2 OFFSET $3", commentsTable)
	err := r.db.Select(&comments, query, postId, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range comments {
		replies, err := r.getRepliesForComment(comments[i].Id)
		if err != nil {
			return nil, err
		}
		comments[i].Replies = replies
	}

	return comments, nil
}

func (r *CommentPostgres) getRepliesForComment(parentCommentID int) ([]models.Comment, error) {
	var replies []models.Comment
	query := fmt.Sprintf("SELECT id, post_id, user_id, parent_comment_id, content "+
		" FROM %s WHERE parent_comment_id = $1", commentsTable)
	err := r.db.Select(&replies, query, parentCommentID)
	if err != nil {
		return nil, err
	}

	for i := range replies {
		subReplies, err := r.getRepliesForComment(replies[i].Id)
		if err != nil {
			return nil, err
		}
		replies[i].Replies = subReplies
	}

	return replies, nil
}
