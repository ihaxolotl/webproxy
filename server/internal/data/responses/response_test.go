package responses

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

const DatabasePath = "./db.sqlite"

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

func testTable() *ResponseTable {
	file, err := os.Create(DatabasePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	db, err := sql.Open("sqlite", DatabasePath)
	if err != nil {
		panic(err)
	}

	table := &ResponseTable{db}
	if err = table.Create(); err != nil {
		panic(err)
	}

	return table
}

func testInsertResponse(table *ResponseTable) (int64, error) {
	return table.Insert(testExampleResponse)
}

func testInsertAndGetResponse(table *ResponseTable) (*Response, error) {
	return table.InsertAndFetch(testExampleResponse)
}

func TestResponseInsert(t *testing.T) {
	table := testTable()
	if _, err := testInsertResponse(table); err != nil {
		t.Fatal(err)
	}
}

func TestResponseInsertAndFetch(t *testing.T) {
	table := testTable()
	if _, err := testInsertAndGetResponse(table); err != nil {
		t.Fatal(err)
	}
}

func TestResponseFetchById(t *testing.T) {
	table := testTable()

	inserted, err := testInsertAndGetResponse(table)
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
