package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "password"
	dbName   = "futbits_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB
var models Models

func TestMain(m *testing.M) {

	// connect to docker; fail if not running
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; check it's running %s", err)
	}

	pool = p

	// set up our docker options
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// get resource
	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {

		var err error

		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}

		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database: %s", err)
	}

	models = New(testDB)

	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	// m -> main.Run() -> run tests from main function
	code := m.Run()

	// clean up docker container
	// For testing comment this if err code block
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables() error {
	// get init.up.sql file
	talbeSQL, err := os.ReadFile("../migrations/init.up.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	// get seed.up.sql file
	seedUpSQL, err := os.ReadFile("../migrations/seed.up.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	// exec init.up.sql file
	_, err = testDB.Exec(string(talbeSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	// exec seed.up.sql file
	_, err = testDB.Exec(string(seedUpSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Test_pingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("can't ping database")
	}
}

func TestPostgresDBRepoInsertGame(t *testing.T) {
	var game Game

	gameBody := Game{
		FieldID:  "828378ed-90f2-453c-af77-7706a25519cb",
		GameDate: time.Date(2024, 2, 16, 24, 0, 0, 0, time.UTC),
		Score:    "12-8",
		GroupID:  "727378ed-20f2-453c-af77-7706a63419cb",
	}

	_, err := game.CreateGame(gameBody)
	if err != nil {
		t.Errorf("Create game returned an error: %s", err)
	}

}

// func CreateGame(testGame services.Game) {
// 	// panic("unimplemented")
// }
