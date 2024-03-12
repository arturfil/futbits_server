package db

import "database/sql"

type DatabaseRepo interface {
    Connect(dsn string) (*DB, error)
    checkDB(d *sql.DB)
}
