package data

import (
	"database/sql"
	"os"

	"github.com/ihaxolotl/webproxy/internal/data/projects"
	"github.com/ihaxolotl/webproxy/internal/data/requests"
	"github.com/ihaxolotl/webproxy/internal/data/responses"
	_ "modernc.org/sqlite"
)

const DatabasePath = "./db.sqlite"

type Table interface {
	Create() error
}

type Database struct {
	conn      *sql.DB
	Projects  *projects.ProjectsTable
	Requests  *requests.RequestsTable
	Responses *responses.ResponseTable
}

func New() *Database {
	return &Database{}
}

// Connect opens the SQLite3 database.
func (db *Database) connect() (*sql.DB, error) {
	file, err := os.Create(DatabasePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return sql.Open("sqlite", DatabasePath)
}

// SetupDatabase connects to the database instance and creates the
// necessary tables.
func (db *Database) Setup() (err error) {
	var tables []Table

	if db.conn, err = db.connect(); err != nil {
		return err
	}

	db.Projects = projects.New(db.conn)
	db.Requests = requests.New(db.conn)
	db.Responses = responses.New(db.conn)

	tables = []Table{db.Projects, db.Requests, db.Responses}
	for _, t := range tables {
		if err = t.Create(); err != nil {
			return err
		}
	}

	return err
}
