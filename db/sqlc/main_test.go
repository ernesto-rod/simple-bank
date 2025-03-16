package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/ernesto-rod/simple-bank/util"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:tr4nsactD3@localhost:5433/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: %w\n", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
