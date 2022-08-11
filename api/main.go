package main

import (
	"chi_soccer/internal/data"
	"chi_soccer/internal/driver"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	port string
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	models   data.Models
}

func main() {
	godotenv.Load(".env")
	var cfg config
	port := os.Getenv("PORT")
	cfg.port = port

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := os.Getenv("DSN")
	db, err := driver.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	defer db.SQL.Close()

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		models:   data.New(db.SQL),
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) serve() error {
	port := os.Getenv("PORT")
	app.infoLog.Println("API listening on port", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	return srv.ListenAndServe()
}
