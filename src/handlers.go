package main

import "net/http"

func (app *application) listFilesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Here are all the files"))
}

func (app *application) displayFileHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("fileId")
	w.Write([]byte("List File: " + id))
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	health := &map[string]string{
		"Status":      "Available",
		"Environment": app.cfg.env,
	}
	if err := app.writeJSON(w, health); err != nil {
		app.logger.Printf("Unable to write JSON to http: %v\n", err)
		return
	}
}
