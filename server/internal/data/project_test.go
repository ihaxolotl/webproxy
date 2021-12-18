package data

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
)

var testExampleProject = &Project{
	ID:          uuid.New().String(),
	Title:       "Test Project",
	Description: "This is a test description",
	Created:     time.Now(),
}

func testInsertProject(db *sql.DB) (int64, error) {
	return InsertProject(db, testExampleProject)
}

func testInsertAndGetProject(db *sql.DB) error {
	return InsertAndGetProject(db, testExampleProject)
}

func TestInsertProject(t *testing.T) {
	db := testDatabase()
	if _, err := testInsertProject(db); err != nil {
		t.Fatal(err)
	}
}

func TestInsertAndGetProject(t *testing.T) {
	db := testDatabase()

	if err := testInsertAndGetProject(db); err != nil {
		t.Fatal(err)
	}
}

func TestGetProjects(t *testing.T) {
	db := testDatabase()
	n := 3

	for i := 0; i < n; i++ {
		if _, err := testInsertProject(db); err != nil {
			t.Fatal(err)
		}
	}

	projects, err := GetProjects(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(projects) != n {
		t.Fatalf("fatal: %d results expected, %d resulted returned.\n", n, len(projects))
	}
}

func TestGetProjectById(t *testing.T) {
	db := testDatabase()

	inserted := testExampleProject
	if err := testInsertAndGetProject(db); err != nil {
		t.Fatal(err)
	}

	fetched, err := GetProjectById(db, inserted.ID)
	if err != nil {
		t.Fatal(err)
	}

	if inserted.ID != fetched.ID {
		t.Fatalf("fatal: inserted ID (%s) does not match fetched ID (%s).\n", inserted.ID, fetched.ID)
	}
}
