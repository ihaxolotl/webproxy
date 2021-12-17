package data

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

var testExampleRequest = &Request{
	ID:        uuid.New().String(),
	ProjectID: uuid.New().String(),
	Method:    "GET",
	Domain:    "localhost",
	IPAddr:    "127.0.0.1",
	Length:    18,
	Edited:    true,
	Timestamp: time.Now(),
	Comment:   "SQL injection.",
	Raw:       "GET / HTTP/1.1\r\n\r\n",
}

func testInsertRequest(db *sql.DB) (int64, error) {
	return InsertRequest(db, testExampleRequest)
}

func testInsertAndGetRequest(db *sql.DB) (*Request, error) {
	return InsertAndGetRequest(db, testExampleRequest)
}

func TestInsertRequest(t *testing.T) {
	db := testDatabase()
	if _, err := testInsertRequest(db); err != nil {
		t.Fatal(err)
	}
}

func TestInsertAndGetRequest(t *testing.T) {
	db := testDatabase()
	if _, err := testInsertAndGetRequest(db); err != nil {
		t.Fatal(err)
	}
}

func TestGetRequestById(t *testing.T) {
	db := testDatabase()

	inserted, err := testInsertAndGetRequest(db)
	if err != nil {
		t.Fatal(err)
	}

	fetched, err := GetRequestById(db, inserted.ID)
	if err != nil {
		t.Fatal(err)
	}

	if inserted.ID != fetched.ID {
		t.Fatalf("fatal: inserted ID (%s) does not match fetched ID (%s).\n", inserted.ID, fetched.ID)
	}

	fmt.Printf("%+#v\n", fetched)
}
