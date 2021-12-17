package data

import (
	"database/sql"
	"testing"
)

func testDatabase() (db *sql.DB) {
	var err error
	if db, err = SetupDatabase(); err != nil {
		panic(err)
	}

	return db
}

func testInsertProject(db *sql.DB) (int64, error) {
	testTitle := "Test Project"
	testDesciption := "This is a test description"

	return InsertProject(db, testTitle, testDesciption)
}

func testInsertAndGetProject(db *sql.DB) (*Project, error) {
	testTitle := "Test Project"
	testDesciption := "This is a test description"

	return InsertAndGetProject(db, testTitle, testDesciption)
}

func TestInsertProject(t *testing.T) {
	db := testDatabase()
	if _, err := testInsertProject(db); err != nil {
		t.Fatal(err)
	}
}

func TestInsertAndGetProject(t *testing.T) {
	db := testDatabase()

	if _, err := testInsertAndGetProject(db); err != nil {
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

	inserted, err := testInsertAndGetProject(db)
	if err != nil {
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
