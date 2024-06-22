package handler

import (
	"TalkBoard/pkg/service"
	"TalkBoard/schema_graphql"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitGraphQL() http.Handler {
	schema_graphql.Initialize()
	schema := h.createSchema()

	return handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})
}

func (h *Handler) createSchema() graphql.Schema {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"getAllPosts": &graphql.Field{
				Type:        graphql.NewList(schema_graphql.PostType),
				Description: "Get all posts",
				Resolve:     h.getAllPosts,
			},
			"getPostById": &graphql.Field{
				Type:        schema_graphql.PostWithCommentsType,
				Description: "Get a post by ID with comments",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: h.getPostById,
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createComment": &graphql.Field{
				Type:        schema_graphql.CommentType,
				Description: "Create a new comment",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(schema_graphql.Comment),
					},
				},
				Resolve: h.createComment,
			},
			"createPost": &graphql.Field{
				Type:        schema_graphql.PostType,
				Description: "Create a new post",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(schema_graphql.Post),
					},
				},
				Resolve: h.createPost,
			},
			"signUp": &graphql.Field{
				Type:        schema_graphql.UserType,
				Description: "Sign up a new user",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(schema_graphql.SignUpInput),
					},
				},
				Resolve: h.signUp,
			},
			"signIn": &graphql.Field{
				Type:        schema_graphql.SignInResponse,
				Description: "Sign in an existing user",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(schema_graphql.SignInInput),
					},
				},
				Resolve: h.signIn,
			},
		},
	})

	schemaConfig := graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		panic(err)
	}

	return schema
}
