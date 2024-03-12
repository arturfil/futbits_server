package handlers

import (
	"chi_soccer/db"
	"chi_soccer/helpers"
	"chi_soccer/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
)

var app services.Application

type Application struct {
	Config services.Config
	Models services.Models
}

func TestMain(m *testing.M) {
   dsn := "host=localhost port=5432 user=root password=secret dbname=chi_soccerdb sslmode=disable timezone=UTC connect_timeout=5" 


	dbConn, err := db.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	var app Application    

    app.Models = services.New(dbConn.DB)

	os.Exit(m.Run())
 
}

func (app *Application) TestServe() error {
	helpers.MessageLogs.InfoLog.Println("API listening on port", "8081")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", "8081"),
		Handler: Routes(),
	}
	return srv.ListenAndServe()
}


