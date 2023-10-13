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

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	cfg    config
	logger *log.Logger
	model  data.Model
}

func main() {

	fmt.Println(os.Getwd())

	var cfg config

	flag.IntVar(&cfg.port, "port", 8000, "Port Number")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://gapi:freeroam@localhost/gofy?sslmode=disable", "DB DataSource Name")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	db, err := OpenDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Println("Database Connection Pool Established")

	model := data.NewModel(db)

	app := &application{
		cfg:    cfg,
		logger: logger,
		model:  model,
	}

	router := app.routes()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	err = srv.ListenAndServe()
	if err != nil {
		app.logger.Fatal(err)
	}
}

func OpenDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, err
}
