package data

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

// CreateProjectTable creates the "projects" table if it doesn't already exist.
func CreateProjectTable(db *sql.DB) (err error) {
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY NOT NULL UNIQUE,
			title TEXT NOT NULL,
			description TEXT NULL,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)

	return err
}

// InsertProject inserts a new project into the database.
func InsertProject(db *sql.DB, p *Project) (rowid int64, err error) {
	var stmt *sql.Stmt

	stmt, err = db.Prepare(`
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

func InsertAndGetProject(db *sql.DB, p *Project) (err error) {
	var (
		stmt  *sql.Stmt
		rowid int64
	)

	rowid, err = InsertProject(db, p)
	if err != nil {
		return err
	}

	stmt, err = db.Prepare(`
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

// GetProjects returns all projects in the database.
// TODO(Brett): Add pagination.
func GetProjects(db *sql.DB) (projects []Project, err error) {
	rows, err := db.Query(`
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

// GetProjectById returns a single project or nil by a string id.
func GetProjectById(db *sql.DB, id string) (project *Project, err error) {
	var (
		stmt *sql.Stmt
		p    *Project
	)

	stmt, err = db.Prepare(`
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
