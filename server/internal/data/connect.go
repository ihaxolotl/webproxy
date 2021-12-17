package data

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const DatabasePath = "./db.sqlite"

// Connect opens the SQLite3 database.
func Connect() (*sql.DB, error) {
	file, err := os.Create(DatabasePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return sql.Open("sqlite", DatabasePath)
}

// SetupDatabase connects to the database instance and creates the
// necessary tables.
func SetupDatabase() (db *sql.DB, err error) {
	db, err = Connect()
	if err != nil {
		return nil, err
	}

	if err = CreateProjectTable(db); err != nil {
		return nil, err
	}

	return db, err
}
