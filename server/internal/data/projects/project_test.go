package projects

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

const DatabasePath = "/tmp/db.sqlite"

var testExampleProject = &Project{
	ID:          uuid.New().String(),
	Title:       "Test Project",
	Description: "This is a test description",
	Created:     time.Now(),
}

func testTable() *ProjectsTable {
	file, err := os.Create(DatabasePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	db, err := sql.Open("sqlite", DatabasePath)
	if err != nil {
		panic(err)
	}

	table := &ProjectsTable{db}
	if err = table.Create(); err != nil {
		panic(err)
	}

	return table
}

func testInsertProject(table *ProjectsTable) (int64, error) {
	return table.Insert(testExampleProject)
}

func testInsertAndGetProject(table *ProjectsTable) error {
	return table.InsertAndFetch(testExampleProject)
}

func TestProjectInsert(t *testing.T) {
	table := testTable()
	if _, err := testInsertProject(table); err != nil {
		t.Fatal(err)
	}
}

func TestProjectInsertAndFetch(t *testing.T) {
	table := testTable()
	if err := testInsertAndGetProject(table); err != nil {
		t.Fatal(err)
	}
}

func TestProjectFetch(t *testing.T) {
	table := testTable()
	n := 3

	for i := 0; i < n; i++ {
		if _, err := testInsertProject(table); err != nil {
			t.Fatal(err)
		}
	}

	projects, err := table.Fetch()
	if err != nil {
		t.Fatal(err)
	}

	if len(projects) != n {
		t.Fatalf("fatal: %d results expected, %d resulted returned.\n", n, len(projects))
	}
}

func TestProjectFetchById(t *testing.T) {
	table := testTable()

	inserted := testExampleProject
	if err := testInsertAndGetProject(table); err != nil {
		t.Fatal(err)
	}

	fetched, err := table.FetchById(inserted.ID)
	if err != nil {
		t.Fatal(err)
	}

	if inserted.ID != fetched.ID {
		t.Fatalf("fatal: inserted ID (%s) does not match fetched ID (%s).\n", inserted.ID, fetched.ID)
	}
}
