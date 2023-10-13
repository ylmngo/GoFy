package main

import (
	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	router := mux.NewRouter()

	router.Use(app.authenticate)

	router.HandleFunc("/", app.displayFilesHandler).Methods("GET")
	router.HandleFunc("/health", app.healthCheckHandler).Methods("GET")
	router.HandleFunc("/register", app.registerHandler).Methods("POST")
	router.HandleFunc("/login", app.loginHandler).Methods("POST")
	router.HandleFunc("/uploads", app.uploadFileHandler).Methods("POST")
	router.HandleFunc("/file/{filename}", app.displayFileHandler).Methods("GET")

	return router
}
