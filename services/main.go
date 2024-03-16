package services

import (
	"database/sql"
	"time"
)

// time for db process any transaction
const dbTimeout = time.Second * 3

var db *sql.DB

// create a new db pool to cache data
func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		User:  User{},
		Field: &Field{},
	}
}

type Config struct {
	Port string
}

type Models struct {
	User          User
	JsonResponse  JsonResponse
	Report        Report
	Field         FieldRepo
	Profile       Profile
	Member        Member
	Group         Group
	TokenResponse TokenResponse
	Game          Game
}
