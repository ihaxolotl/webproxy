package requests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

const DatabasePath = "/tmp/db.sqlite"

var testExampleRequest = &Request{
	ID:         uuid.New().String(),
	ProjectID:  uuid.New().String(),
	ResponseID: uuid.New().String(),
	Method:     "GET",
	Domain:     "localhost",
	IPAddr:     "127.0.0.1",
	URL:        "/",
	Length:     18,
	Edited:     true,
	Timestamp:  time.Now(),
	Comment:    "SQL injection.",
	Raw:        "GET / HTTP/1.1\r\n\r\n",
}

func testTable() *RequestsTable {
	file, err := os.Create(DatabasePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	db, err := sql.Open("sqlite", DatabasePath)
	if err != nil {
		panic(err)
	}

	table := &RequestsTable{db}
	if err = table.Create(); err != nil {
		panic(err)
	}

	return table
}

func testInsertRequest(table *RequestsTable) (int64, error) {
	return table.Insert(testExampleRequest)
}

func testInsertAndGetRequest(table *RequestsTable) (*Request, error) {
	return table.InsertAndFetch(testExampleRequest)
}

func TestRequestInsert(t *testing.T) {
	table := testTable()
	if _, err := testInsertRequest(table); err != nil {
		t.Fatal(err)
	}
}

func TestRequestInsertAndFetch(t *testing.T) {
	table := testTable()
	if _, err := testInsertAndGetRequest(table); err != nil {
		t.Fatal(err)
	}
}

func TestRequestFetchById(t *testing.T) {
	table := testTable()

	inserted, err := testInsertAndGetRequest(table)
	if err != nil {
		t.Fatal(err)
	}

	fetched, err := table.FetchById(inserted.ID)
	if err != nil {
		t.Fatal(err)
	}

	if inserted.ID != fetched.ID {
		t.Fatalf("fatal: inserted ID (%s) does not match fetched ID (%s).\n", inserted.ID, fetched.ID)
	}

	fmt.Printf("%+#v\n", fetched)
}
