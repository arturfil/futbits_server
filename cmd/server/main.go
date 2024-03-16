package main

import (
	"chi_soccer/db"
	"chi_soccer/models"
	"chi_soccer/services"
	"log"
	"os"
)

func main() {
	var cfg services.Config
	var db db.DB
	port := os.Getenv("PORT")
	cfg.Port = port

	dsn := os.Getenv("DSN")
	dbConn, err := db.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	defer dbConn.DB.Close()

	var app = &models.Application{
		Config: cfg,
		Models: services.New(dbConn.DB),
	}

	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
