package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, data any) error {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}
