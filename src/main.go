package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	env  string
	dsn  string
	port int
}

type application struct {
	cfg    config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8000, "Service Port Number")
	flag.StringVar(&cfg.env, "env", "Development", "Environment (Development|Production)")
	flag.StringVar(&cfg.dsn, "dsn", "postgres://gofy:freeroam@localhost/gofy?sslmode=disable", "DB datasource name")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		cfg:    cfg,
		logger: logger,
	}

	router := app.routes()

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	if err := srv.ListenAndServe(); err != nil {
		app.logger.Fatal(err)
	}
}
