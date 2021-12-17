package data

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Project respresents the data for an engagement.
type Project struct {
	ID          string    // Unique ID of the project.
	Title       string    // Title of the project.
	Description string    // Description of the project.
	Created     time.Time // Timestamp for when the project was created.
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
func InsertProject(db *sql.DB, title, description string) (rowid int64, err error) {
	stmt, err := db.Prepare(`
		INSERT INTO projects(
			id, title, description, created
		) VALUES (?, ?, ?, ?);
	`)

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(uuid.New().String(), title, description, time.Now())
	return res.LastInsertId()
}

func InsertAndGetProject(db *sql.DB, title, description string) (project *Project, err error) {
	var (
		stmt  *sql.Stmt
		rowid int64
		p     *Project
	)

	rowid, err = InsertProject(db, title, description)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	p = &Project{}
	if err = stmt.QueryRow(rowid).Scan(
		&p.ID,
		&p.Title,
		&p.Description,
		&p.Created,
	); err != nil {
		return nil, err
	}

	return p, err
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

		if err = rows.Scan(&p.ID, &p.Title, &p.Description, &p.Created); err != nil {
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
	if err = stmt.QueryRow(id).Scan(&p.ID, &p.Title, &p.Description, &p.Created); err != nil {
		return nil, err
	}

	return p, err
}
