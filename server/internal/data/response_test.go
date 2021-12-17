package data

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

var testExampleResponse = &Response{
	ID:        uuid.New().String(),
	ProjectID: uuid.New().String(),
	Status:    200,
	Length:    18,
	Elapsed:   100,
	Edited:    true,
	Timestamp: time.Now(),
	Mimetype:  "text/html",
	Comment:   "SQL injection.",
	Raw:       "HTTP/1.0 200 OK\r\nServer: SimpleHTTP/0.6 Python/3.10.1\r\nDate: Fri, 17 Dec 2021 10:45:06 GMT\r\nContent-type: text/html; charset=utf-8\r\nContent-Length: 0\r\n\r\n",
}

func testInsertResponse(db *sql.DB) (int64, error) {
	return InsertResponse(db, testExampleResponse)
}

func testInsertAndGetResponse(db *sql.DB) (*Response, error) {
	return InsertAndGetResponse(db, testExampleResponse)
}

func TestInsertResponse(t *testing.T) {
	db := testDatabase()
	if _, err := testInsertResponse(db); err != nil {
		t.Fatal(err)
	}
}

func TestInsertAndGetResponse(t *testing.T) {
	db := testDatabase()
	if _, err := testInsertAndGetResponse(db); err != nil {
		t.Fatal(err)
	}
}

func TestGetResponseById(t *testing.T) {
	db := testDatabase()

	inserted, err := testInsertAndGetResponse(db)
	if err != nil {
		t.Fatal(err)
	}

	fetched, err := GetResponseById(db, inserted.ID)
	if err != nil {
		t.Fatal(err)
	}

	if inserted.ID != fetched.ID {
		t.Fatalf("fatal: inserted ID (%s) does not match fetched ID (%s).\n", inserted.ID, fetched.ID)
	}

	fmt.Printf("%+#v\n", fetched)
}
