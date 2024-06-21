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
			// Простой пустой запрос, чтобы избежать ошибки
			"dummyQuery": &graphql.Field{
				Type:        graphql.String,
				Description: "A dummy query",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "Hello, World!", nil
				},
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createPost": &graphql.Field{
				Type:        schema_graphql.PostType,
				Description: "Create a new post",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(schema_graphql.CreatePostInput),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return h.createPost(p)
				},
			},
			"signUp": &graphql.Field{
				Type:        schema_graphql.UserType, // Возвращаемый тип данных (пользователь)
				Description: "Sign up a new user",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(schema_graphql.SignUpInput),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return h.signUp(p)
				},
			},
			"signIn": &graphql.Field{
				Type:        schema_graphql.SignInResponse, // Возвращаемый тип данных (токен)
				Description: "Sign in an existing user",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(schema_graphql.SignInInput),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return h.signIn(p)
				},
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
