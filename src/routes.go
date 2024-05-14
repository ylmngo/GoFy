package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/list", app.listFilesHandler)
	router.HandleFunc("/list/{filename}", app.displayFileHandler)
	router.HandleFunc("/health", app.healthCheckHandler)
	router.HandleFunc("/upload", app.uploadFileHandler)

	return router
}
