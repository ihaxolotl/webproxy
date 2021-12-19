package projects

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Project respresents the data for an engagement.
type Project struct {
	ID          string    `json:"id"`         // Unique ID of the project.
	Title       string    `json:"title"`      // Title of the project.
	Description string    `json:"decription"` // Description of the project.
	Created     time.Time `json:"created"`    // Timestamp for when the project was created.
}

type ProjectsTable struct {
	db *sql.DB
}

func New(db *sql.DB) *ProjectsTable {
	return &ProjectsTable{db}
}

// Create creates the "projects" table if it doesn't already exist.
func (t ProjectsTable) Create() (err error) {
	_, err = t.db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY NOT NULL UNIQUE,
			title TEXT NOT NULL,
			description TEXT NULL,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)

	return err
}

// Insert inserts a new record into the projects table and returns the last inserted
// rowid or an error.
func (t ProjectsTable) Insert(p *Project) (rowid int64, err error) {
	var stmt *sql.Stmt

	stmt, err = t.db.Prepare(`
		INSERT INTO projects(
			id, title, description, created
		) VALUES (?, ?, ?, ?);
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	p.ID = uuid.New().String()
	p.Created = time.Now()

	res, err := stmt.Exec(
		p.ID,
		p.Title,
		p.Description,
		p.Created,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// InsertAndFetch inserts a new record into the projects table and returns the full
// record or an error.
func (t ProjectsTable) InsertAndFetch(p *Project) (err error) {
	var (
		stmt  *sql.Stmt
		rowid int64
	)

	rowid, err = t.Insert(p)
	if err != nil {
		return err
	}

	stmt, err = t.db.Prepare(`
		SELECT
			id, title, description, created
		FROM
			projects
		WHERE
			rowid = ?
		LIMIT 0, 1;
	`)
	if err != nil {
		return err
	}

	return stmt.QueryRow(rowid).Scan(
		&p.ID,
		&p.Title,
		&p.Description,
		&p.Created,
	)
}

// Fetch returns all records from the projects table.
// TODO(Brett): Add pagination.
func (t ProjectsTable) Fetch() (projects []Project, err error) {
	var rows *sql.Rows

	rows, err = t.db.Query(`
		SELECT
			id, title, description, created
		FROM
			projects;
	`)

	projects = make([]Project, 0)

	for rows.Next() {
		var p Project

		if err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.Created,
		); err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	return projects, err
}

// FetchById returns a single project or nil by a string id.
func (t ProjectsTable) FetchById(id string) (p *Project, err error) {
	var (
		stmt *sql.Stmt
	)

	stmt, err = t.db.Prepare(`
		SELECT
			id, title, description, created
		FROM 
			projects
		WHERE
			id = ?
		LIMIT 0, 1;
	`)
	if err != nil {
		return nil, err
	}

	p = &Project{}
	if err = stmt.QueryRow(id).Scan(
		&p.ID,
		&p.Title,
		&p.Description,
		&p.Created,
	); err != nil {
		return nil, err
	}

	return p, err
}

func (t ProjectsTable) FetchHistory(id string) (err error) {
	return nil
}
