package schema_graphql

import "github.com/graphql-go/graphql"

var CommentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"userId": &graphql.Field{
				Type: graphql.Int,
			},
			"content": &graphql.Field{
				Type: graphql.String,
			},
			"postId": &graphql.Field{
				Type: graphql.Int,
			},
			"parentCommentId": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

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
