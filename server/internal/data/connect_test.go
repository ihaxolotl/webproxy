package data

import (
	"database/sql"
	"testing"
)

func TestConnect(t *testing.T) {
	if _, err := Connect(); err != nil {
		t.Fatal(err)
	}
}

func TestSetupDatabase(t *testing.T) {
	if _, err := SetupDatabase(); err != nil {
		t.Fatal(err)
	}
}

func testDatabase() (db *sql.DB) {
	var err error
	if db, err = SetupDatabase(); err != nil {
		panic(err)
	}

	return db
}
