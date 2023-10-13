package main

import (
	"errors"
	"gofy/internal/data"
	"gofy/internal/validator"
	"net/http"
	"time"
)

func (app *application) registerHandler(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	err := app.readJSON(writer, request, &input)
	if err != nil {
		app.logger.Fatal("Could not read JSON input")
		return
	}

	user := &data.User{
		Username: input.Username,
		Email:    input.Email,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.logger.Fatal("Could not Hash Password")
		return
	}

	v := validator.New()

	data.ValidateUser(v, user)

	if !v.Valid() {

		for key, val := range v.Errors {
			app.logger.Println(key, val)
		}

		app.logger.Fatal("Not a Valid User")
		return
	}

	err = app.model.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			app.logger.Fatal("Email must be unique")
			return
		default:
			app.logger.Fatal("User could not be inserted")
			return
		}
	}

	err = app.writeJSON(writer, map[string]interface{}{"user": user})
	if err != nil {
		app.logger.Fatal("Server could not respond to the request")
	}
}

func (app *application) loginHandler(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(writer, request, &input)
	if err != nil {
		app.logger.Fatal("Could not read JSON data")
		return
	}

	user, err := app.model.Users.GetByEmail(input.Email)
	if err != nil {
		app.logger.Fatal("Could not retreive User")
		return
	}

	ok, _ := user.Password.Matches(input.Password)
	if !ok {
		app.logger.Fatal("Incorrect Password")
		return
	}

	token, err := app.model.Tokens.GenerateToken(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.logger.Fatal("Error while generating token")
		return
	}

	if token == nil {
		app.logger.Fatal("Token is Nil")
		return
	}

	err = app.model.Tokens.Insert(token)
	if err != nil {
		app.logger.Fatal("Error while inserting Token")
		return
	}

	data := map[string]interface{}{
		"authentication_token": token,
	}

	app.writeJSON(writer, data)
}
