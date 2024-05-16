package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()
	auth := http.NewServeMux()

	auth.HandleFunc("/list/{fileId}", app.displayFileHandler)
	auth.HandleFunc("/list", app.listFilesHandler)
	auth.HandleFunc("/upload", app.uploadFileHandler)

	router.Handle("/", app.authenticate(auth))
	router.HandleFunc("/health", app.healthCheckHandler)
	router.HandleFunc("/register", app.registerUserHandler)
	router.HandleFunc("/login", app.loginHandler)

	nm := app.nestMiddlewares(
		app.recover,
		app.rateLimit,
	)

	return nm(router)
}
