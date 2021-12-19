package history

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ihaxolotl/webproxy/internal/data/projects"
	"github.com/ihaxolotl/webproxy/internal/data/requests"
	"github.com/ihaxolotl/webproxy/internal/data/responses"
	_ "modernc.org/sqlite"
)

const DatabasePath = "/tmp/db.sqlite"

type testTable interface {
	Create() error
}

var testExampleProject = &projects.Project{
	ID:          uuid.New().String(),
	Title:       "Test Project",
	Description: "This is a test description",
	Created:     time.Now(),
}

func testView() *HistoryView {
	file, err := os.Create(DatabasePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	db, err := sql.Open("sqlite", DatabasePath)
	if err != nil {
		panic(err)
	}

	testCreateRequirements(db)

	table := &HistoryView{db}
	if err = table.Create(); err != nil {
		panic(err)
	}

	return table
}

func testCreateRequirements(db *sql.DB) {
	proj := projects.New(db)
	req := requests.New(db)
	res := responses.New(db)
	n := 3

	tables := []testTable{proj, req, res}

	for _, t := range tables {
		if err := t.Create(); err != nil {
			panic(err)
		}
	}

	if err := proj.InsertAndFetch(testExampleProject); err != nil {
		panic(err)
	}

	for i := 0; i < n; i++ {
		reqid := uuid.New().String()
		resid := uuid.New().String()

		testExampleRequest := &requests.Request{
			ID:         reqid,
			ProjectID:  testExampleProject.ID,
			ResponseID: resid,
			Method:     http.MethodGet,
			Domain:     "localhost",
			IPAddr:     "127.0.0.1",
			URL:        "/",
			Length:     18,
			Edited:     true,
			Timestamp:  time.Now(),
			Comment:    "SQL injection.",
			Raw:        "GET / HTTP/1.1\r\n\r\n",
		}

		testExampleResponse := &responses.Response{
			ID:        resid,
			ProjectID: testExampleProject.ID,
			RequestID: reqid,
			Status:    http.StatusOK,
			Length:    18,
			Elapsed:   100,
			Edited:    true,
			Timestamp: time.Now(),
			Mimetype:  "text/html",
			Comment:   "SQL injection.",
			Raw:       "HTTP/1.0 200 OK\r\nServer: SimpleHTTP/0.6 Python/3.10.1\r\nDate: Fri, 17 Dec 2021 10:45:06 GMT\r\nContent-type: text/html; charset=utf-8\r\nContent-Length: 0\r\n\r\n",
		}

		if _, err := req.InsertAndFetch(testExampleRequest); err != nil {
			panic(err)
		}

		if _, err := res.InsertAndFetch(testExampleResponse); err != nil {
			panic(err)
		}
	}
}

func TestHistoryFetch(t *testing.T) {
	view := testView()
	n := 3

	records, err := view.Fetch(testExampleProject.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(records) != n {
		t.Fatalf("fatal: %d results expected, %d results returned.\n", n, len(records))
	}

	fmt.Printf("%v\n", records)
}
