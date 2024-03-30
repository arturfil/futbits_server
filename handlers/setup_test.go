package handlers

import (
	"chi_soccer/helpers"
	"chi_soccer/services"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"testing"
)

type Application struct {
	Config services.Config
	DB     services.DB
}

func TestMain(m *testing.M) {
	// dsn := "host=localhost port=5432 user=root password=secret dbname=chi_soccerdb sslmode=disable timezone=UTC connect_timeout=5"
	// db := services.DB{}

	// dbConn, err := db.ConnectToDB(dsn)
	// if err != nil {
	// 	log.Fatal("Cannot connect to database", err)
	// }

	var app Application
	app.DB = services.DB{
		Models: services.Models{Field: &services.MockField{}},
	}

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

func New(dbPool *sql.DB) services.Models {
	return services.Models{
		Field: &services.MockField{},
	}
}
