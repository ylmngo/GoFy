package main

import (
	"gofy/internal/storage"
	"net/http"
	"os"
)

func (app *application) listFilesHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir("uploads/")
	if err != nil {
		app.logger.Printf("Unable to read entries of uploads dir: %v\n", err)
		return
	}

	files := make([]string, 0)
	for _, entry := range entries {
		files = append(files, entry.Name())
	}

	if err := app.writeJSON(w, http.StatusOK, files, nil); err != nil {
		app.logger.Printf("Unable to write filenames to json: %v\n", err)
		return
	}
}

func (app *application) displayFileHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("filename")

	data, err := storage.Read("uploads/" + name)
	if err != nil {
		app.logger.Printf("Unable to read file from disk: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Write(data)
}

func (app *application) uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		app.logger.Printf("Unable to parse multipart form: %v\n", err)
		return
	}

	mpf, hdr, err := r.FormFile("file")
	if err != nil {
		app.logger.Printf("Unable to get form file: %v\n", err)
		app.writeJSON(w, http.StatusNotFound, "No file given", nil)
		return
	}
	defer mpf.Close()

	if err := storage.Write("uploads/"+hdr.Filename, mpf); err != nil {
		app.logger.Printf("Unable to write file to disk: %v\n", err)
		app.writeJSON(w, http.StatusInternalServerError, "Try again", nil)
		return
	}

	app.writeJSON(w, http.StatusOK, "File Succesfully Uploaded!", nil)
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	health := &map[string]string{
		"Status":      "Available",
		"Environment": app.cfg.env,
	}
	if err := app.writeJSON(w, http.StatusOK, health, nil); err != nil {
		app.logger.Printf("Unable to write JSON to http: %v\n", err)
		return
	}
}
