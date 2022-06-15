package main

import (
	"chi_soccer/internal/data"
	"chi_soccer/internal/driver"
	"chi_soccer/util"
	"fmt"
	"log"
	"net/http"
	"os"
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
	var cfg config
	port_cfg, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("couldn't load config file")
	}
	cfg.port = port_cfg.ServerAddress

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
	app.infoLog.Println("API listening on port", app.config.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.config.port),
		Handler: app.routes(),
	}

	return srv.ListenAndServe()
}
