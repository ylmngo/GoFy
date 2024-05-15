package main

import (
	"fmt"
	"net/http"
)

type Middleware func(next http.Handler) http.Handler

func (app *application) recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				er := fmt.Sprintf("%s", err)
				app.writeJSON(w, http.StatusInternalServerError, er, nil)
				app.logger.Printf("Internal Server Error\n")
				w.Header().Set("Connection", "close")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
