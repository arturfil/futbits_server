package main

import (
	"chi_soccer/db"
	"chi_soccer/services"
	"log"
	"os"
	"testing"
)

var app services.Application

func TestMain(m *testing.M) {

	dsn := os.Getenv("DSN")

	dbConn, err := db.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	var cfg services.Config
	
	var app = &Application{
		Config: cfg,
		Models: services.New(dbConn.DB),
	}
    
	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())

}
