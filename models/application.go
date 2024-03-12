package models

import (
	"chi_soccer/handlers"
	"chi_soccer/helpers"
	"chi_soccer/services"
	"fmt"
	"net/http"
	"os"
)

type Application struct {
	Config services.Config
	Models services.Models
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
