package schema_graphql

import "github.com/graphql-go/graphql"

var PostType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Post",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"content": &graphql.Field{
				Type: graphql.String,
			},
			"accessToComments": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)

var CreatePostInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "CreatePostInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"userId": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"title": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"content": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"accessToComments": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
		},
	},
)
