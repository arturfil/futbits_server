package main

import (
	"chi_soccer/db"
	"chi_soccer/handlers"
	"chi_soccer/helpers"
	"chi_soccer/services"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Config services.Config
	Models services.Models
}

var srv http.Handler

func main() {
	var cfg services.Config
	port := os.Getenv("PORT")
	cfg.Port = port

	dsn := os.Getenv("DSN")
	dbConn, err := db.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	defer dbConn.DB.Close()

	var app = &Application{
		Config: cfg,
		Models: services.New(dbConn.DB),
	}

	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *Application) Serve() error {
	port := os.Getenv("PORT")
	helpers.MessageLogs.InfoLog.Println("API listening on port", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handlers.Routes(),
	}
	return srv.ListenAndServe()
}


