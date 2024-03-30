package mainmain

import (
	"chi_soccer/models"
	"chi_soccer/services"
	"log"
	"os"
)

func main() {

	var cfg services.Config
	var db services.DB

	port := os.Getenv("PORT")
	cfg.Port = port

	dsn := os.Getenv("DSN")
	dbConn, err := db.ConnectToDB(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	defer dbConn.DB.Close()

	var app = &models.Application{
		DB: services.DB{
			Models: services.New(dbConn.DB),
		},
	}

	err = app.Serve()
	if err != nil {
		log.Panic(err)
	}
}
