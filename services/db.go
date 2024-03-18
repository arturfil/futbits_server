package services

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

type DatabaseRepo interface {
	ConnectPosgres(dsn string) (*DB, error)
	checkDB(d *sql.DB)
}

type DB struct {
	DB *sql.DB
    Models Models
}


const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute

func (d *DB) ConnectPostgres(dsn string) (*DB, error) {
    dbConn := &DB{}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifeTime)

	err = checkDB(db)
	if err != nil {
		return nil, err
	}
	dbConn.DB = db
	return dbConn, nil
}

func checkDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		fmt.Println("Error", err)
		return err
	}
	fmt.Printf("\n*** Pinged database successfully! ***\n")
	return nil
}
