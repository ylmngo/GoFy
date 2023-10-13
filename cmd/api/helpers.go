package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) writeJSON(writer http.ResponseWriter, data interface{}) error {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		return err
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonResponse)

	return nil
}

func (app *application) readJSON(writer http.ResponseWriter, request *http.Request, dst interface{}) error {
	err := json.NewDecoder(request.Body).Decode(dst)
	if err != nil {
		return err
	}

	return nil
}
