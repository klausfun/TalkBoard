package handler

import (
	"TalkBoard/pkg/service"
	"context"
	"errors"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtxKey          = "userId"
)

type Middleware struct {
	services *service.Service
}

func NewMiddleware(services *service.Service) *Middleware {
	return &Middleware{services: services}
}

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			http.Error(w, "empty auth header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		if headerParts[0] != "Bearer" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		if headerParts[1] == "" {
			http.Error(w, "token is empty", http.StatusUnauthorized)
			return
		}

		userId, err := m.services.Authorization.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserId(ctx context.Context) (int, error) {
	userId, ok := ctx.Value(userCtxKey).(int)
	if !ok {
		return 0, errors.New("user id not found")
	}
	return userId, nil
}
