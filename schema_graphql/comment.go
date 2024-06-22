package schema_graphql

import (
	"github.com/graphql-go/graphql"
)

var (
	CommentType          *graphql.Object
	PostWithCommentsType *graphql.Object
)

func Initialize() {
	CommentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Comment",
			Fields: graphql.Fields{},
		},
	)

	PostWithCommentsType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "PostWithComments",
			Fields: graphql.Fields{},
		},
	)

	CommentType.AddFieldConfig("id", &graphql.Field{Type: graphql.Int})
	CommentType.AddFieldConfig("userId", &graphql.Field{Type: graphql.Int})
	CommentType.AddFieldConfig("content", &graphql.Field{Type: graphql.String})
	CommentType.AddFieldConfig("postId", &graphql.Field{Type: graphql.Int})
	CommentType.AddFieldConfig("parentCommentId", &graphql.Field{Type: graphql.Int})
	CommentType.AddFieldConfig("replies", &graphql.Field{Type: graphql.NewList(CommentType)})

	PostWithCommentsType.AddFieldConfig("post", &graphql.Field{Type: PostType})
	PostWithCommentsType.AddFieldConfig("comments", &graphql.Field{Type: graphql.NewList(CommentType)})
}

var Comment = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "CommentInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"parentCommentId": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"postId": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"userId": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"content": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)
