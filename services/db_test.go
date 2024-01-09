package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

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
    opts := dockertest.RunOptions {
        Repository: "postgres",
        Tag: "14.5",
        Env: []string {
            "POSTGRES_USER=" + user,
            "POSTGRES_PASSWORD=" + password,
            "POSTGRES_DB=" + dbName,
        },
        ExposedPorts: []string{"5432"},
        PortBindings: map[docker.Port][]docker.PortBinding {
            "5432": {
                {HostIP: "0.0.0.0", HostPort: port},
            },
        },
    }

    // get resource
    resource, err = pool.RunWithOptions(&opts)
    if err != nil {
        // _ = pool.Purge(resource)
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

     // clean up
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
    var fields = models.Field

    // var game services.Game
    // dbConn, err := testDB
    // services.New()

    // var fields services.Field

    // ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
    // defer cancel()

    // query := `select * from fields`
    // allFields, err := testDB.QueryContext(ctx, query)
    // if err != nil {
    //     t.Errorf("create game returned error: %s", err)
    // }

    //:::::::::::::::::::::::::: HERE :::::::::::::::::::::::::::::
     allFields, err := fields.GetAllFields()
     if err != nil {
         t.Errorf("create game returned error: %s", err)
     }                                                      // fmt.Println(allFields)

    fmt.Println("all fields ->", allFields)

    // testGame := services.Game{
    //     FieldID  : "60950ecc-b886-4428-bf82-4503b8349175",
    //     GameDate : time.Date(2023, 12, 15, 7, 0, 0, 0, time.UTC),
    //     Score    : "10-8",
    //     GroupID  : "d676a368-95ee-41ff-a884-a68cc708de64",
    // }

    // fmt.Println(testGame)

    // gameResp, err := game.CreateGame(testGame)
    // if err != nil {
    //     t.Errorf("create game returned error: %s", err)
    // }

    //resp, err := game.GetAllGames("35069b00-a556-4b4b-acb2-d842007b8ffa")
    //if err != nil {
    //    t.Errorf("create game returned error: %s", err)
    //}

    // fmt.Println(gameResp)

}

//func CreateGame(testGame services.Game) {
//	panic("unimplemented")
//}
