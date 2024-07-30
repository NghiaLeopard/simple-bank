package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDrive = "postgres";
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQuery *Queries
var testDB *sql.DB
 
func TestMain(m *testing.M) {
	var err error
	testDB,err = sql.Open(dbDrive,dbSource)

	if err != nil {
		log.Fatal("Cannnot connect to db",err)
	}

	testQuery = New(testDB)

	os.Exit(m.Run())
}