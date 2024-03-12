package db

import "database/sql"

type DBMock struct {
    DB *sql.DB
}

func (m *DBMock) ConnectToPostgres(dsn string) (*DB, error) {
    return nil, nil
}
