package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"gofy/internal/data"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
	env  string
	dsn  string
	port int
}

type application struct {
	cfg    config
	logger *log.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8000, "Service Port Number")
	flag.StringVar(&cfg.env, "env", "Development", "Environment (Development|Production)")
	flag.StringVar(&cfg.dsn, "dsn", "postgres://gofy:freeroam@localhost/gofy?sslmode=disable", "DB datasource name")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	logger.Println("database connection pool established")

	app := &application{
		cfg:    cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	router := app.routes()

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.Printf("Starting development server on :%d", app.cfg.port)

	if err := srv.ListenAndServe(); err != nil {
		app.logger.Fatal(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	return db, err
}
