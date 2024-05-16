package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/list", app.listFilesHandler)
	router.HandleFunc("/list/{fileId}", app.displayFileHandler)
	router.HandleFunc("/health", app.healthCheckHandler)
	router.HandleFunc("/upload", app.uploadFileHandler)

	router.HandleFunc("/register", app.registerUserHandler)

	nm := app.nestMiddlewares(
		app.recover,
		app.rateLimit,
	)

	return nm(router)
}
