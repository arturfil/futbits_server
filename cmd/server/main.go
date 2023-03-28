package main

import (
	"chi_soccer/db"
	"chi_soccer/helpers"
	"chi_soccer/routes"
	"chi_soccer/services"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Port string
}

type Application struct {
	Config Config
	Models services.Models
}

func (app *Application) Serve() error {
	port := os.Getenv("PORT")
	helpers.MessageLogs.InfoLog.Println("API listening on port", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: routes.Routes(),
	}
	return srv.ListenAndServe()
}

func main() {
	var cfg Config
	port := os.Getenv("PORT")
	cfg.Port = port

	dsn := os.Getenv("DSN")
	dbConn, err := db.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	defer dbConn.DB.Close()

	app := &Application{
		Config: cfg,
		Models: services.New(dbConn.DB),
	}

	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
