package main

import (
	"errors"
	"gofy/internal/data"
	"gofy/internal/validator"
	"net/http"
	"strings"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = app.setContextKey(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.Header().Set("WWW-Authenticate", "Bearer")
			app.logger.Fatal("Invalid Authentication Response")
			return
		}

		token := headerParts[1]

		v := validator.New()
		if data.ValidatePlainTextToken(v, token); !v.Valid() {
			w.Header().Set("WWW-Authenticate", "Bearer")
			app.logger.Fatal("Invalid Authentication Token Response")
			return
		}

		user, err := app.model.Users.GetByToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				w.Header().Set("WWW-Authenticate", "Bearer")
				app.logger.Fatal("Invalid Authentication Response")
			default:
				app.logger.Fatal("Error while retrieving User from token")
			}
			return
		}

		r = app.setContextKey(r, user)
		next.ServeHTTP(w, r)
	})
}
