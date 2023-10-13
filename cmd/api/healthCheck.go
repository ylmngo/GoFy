package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(writer http.ResponseWriter, request *http.Request) {
	user := app.getContextUser(request)
	healthData := &map[string]string{
		"Status":      "Available",
		"User":        user.Username,
		"Environment": app.cfg.env,
		"Version":     version,
	}

	err := app.writeJSON(writer, healthData)
	if err != nil {
		app.logger.Fatal("Could not format data to JSON")
	}
}
