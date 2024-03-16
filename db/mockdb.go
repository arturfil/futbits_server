package db

import (
	"database/sql"
	"fmt"
)

type DBMock struct {
	DB *sql.DB
}

func (m *DBMock) ConnectPostgres(dsn string) (*DB, error) {
	return nil, nil
}

func (m *DBMock) checkDB() error {
	fmt.Println("*** Test db ping worked")
	return nil
}
