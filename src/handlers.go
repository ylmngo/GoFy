package main

import (
	"gofy/internal/data"
	"gofy/internal/storage"
	"net/http"
	"strconv"
)

func (app *application) listFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := app.models.Files.GetAll()
	if err != nil {
		app.logger.Printf("Unable to get files from database: %v\n", err)
		app.writeJSON(w, http.StatusInternalServerError, "Please Try again later", nil)
		return
	}

	app.writeJSON(w, http.StatusOK, files, nil)
}

func (app *application) displayFileHandler(w http.ResponseWriter, r *http.Request) {
	v := r.PathValue("fileId")
	id, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, "Invalid Id", nil)
		app.logger.Printf("Unable to parse path value to id: %v\n", err)
		return
	}

	f, err := app.models.Files.Get(id)
	if err != nil {
		app.writeJSON(w, http.StatusNotFound, "Invalid Id", nil)
		app.logger.Printf("Unable to get file by id = %d: %v\n", id, err)
		return
	}

	data, err := storage.Read("uploads/" + f.Name)
	if err != nil {
		app.logger.Printf("Unable to read file from disk: %v\n", err)
		return
	}

	ct := http.DetectContentType(data)
	w.Header().Set("Content-Type", ct)
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

	f := &data.File{
		Name:     hdr.Filename,
		Metadata: "Something about the file",
	}

	if err := app.models.Files.Insert(f, 1); err != nil {
		app.logger.Printf("Unable to insert file to database: %v\n", err)
		app.writeJSON(w, http.StatusInternalServerError, "Failed to insert to database", nil)
		return
	}

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
