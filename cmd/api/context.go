package main

import (
	"context"
	"gofy/internal/data"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) setContextKey(request *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(request.Context(), userContextKey, user)
	return request.WithContext(ctx)
}

func (app *application) getContextUser(request *http.Request) *data.User {
	user, ok := request.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
