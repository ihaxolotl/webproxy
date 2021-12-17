package data

import (
	"database/sql"
	"time"
)

// Response represents an HTTP response and its metadata that has
// been intercepted by the proxy.
type Response struct {
	ID        string    // Unique ID of the response.
	ProjectID string    // Unique ID of the parent project.
	Status    int16     // HTTP status code of the response.
	Length    int64     // Length of the response in bytes.
	Elapsed   int64     // Time elapsed since request was sent until response.
	Edited    bool      // Flag for whether the response was modified or not.
	Timestamp time.Time // Time the response was received.
	Mimetype  string    // Mime-type of the response body data.
	Comment   string    // User-supplied comment on the response.
	Raw       string    // Raw response bytes.
}

// CreateResponseTable creates the "responses" table if it doesn't already exist.
func CreateResponseTable(db *sql.DB) (err error) {
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS responses (
			id TEXT PRIMARY KEY NOT NULL UNIQUE,
			projectid TEXT NOT NULL,
			status INTEGER NOT NULL,
			length INTEGER NOT NULL,
			elapsed INTEGER NOT NULL,
			edited BOOLEAN NOT NULL CHECK (edited IN (0, 1)),
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			comment TEXT,
			raw TEXT NOT NULL
		);
	`)

	return err
}

// InsertResponse inserts a new entry into the response table and returns
// the last inserted rowid and an error.
func InsertResponse(db *sql.DB, resp *Response) (rowid int64, err error) {
	var (
		stmt *sql.Stmt
		res  sql.Result
	)

	stmt, err = db.Prepare(`
		INSERT INTO responses(
			id,
			projectid,
			status,
			length,
			elapsed,
			edited,
			timestamp,
			comment,
			raw
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?
		);
	`)
	if err != nil {
		return 0, err
	}

	res, err = stmt.Exec(
		resp.ID,
		resp.ProjectID,
		resp.Status,
		resp.Length,
		resp.Elapsed,
		resp.Edited,
		resp.Timestamp,
		resp.Comment,
		resp.Raw,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func InsertAndGetResponse(db *sql.DB, r *Response) (resp *Response, err error) {
	var (
		stmt  *sql.Stmt
		rowid int64
	)

	rowid, err = InsertResponse(db, r)
	if err != nil {
		return nil, err
	}

	stmt, err = db.Prepare(`
		SELECT 
			id,
			projectid,
			status,
			length,
			elapsed,
			edited,
			timestamp,
			comment,
			raw
		FROM
			responses
		WHERE
			rowid = ?
		LIMIT 0, 1;

	`)
	if err != nil {
		return nil, err
	}

	resp = &Response{}
	err = stmt.QueryRow(rowid).Scan(
		&resp.ID,
		&resp.ProjectID,
		&resp.Status,
		&resp.Length,
		&resp.Elapsed,
		&resp.Edited,
		&resp.Timestamp,
		&resp.Comment,
		&resp.Raw,
	)

	return resp, err
}

// GetResponseById queries and returns the entry from the responses table matching
// the supplied id parameter and optionally, an error.
func GetResponseById(db *sql.DB, id string) (resp *Response, err error) {
	var stmt *sql.Stmt

	stmt, err = db.Prepare(`
		SELECT 
			id,
			projectid,
			status,
			length,
			elapsed,
			edited,
			timestamp,
			comment,
			raw
		FROM
			responses	
		WHERE
			id = ?
		LIMIT 0, 1;
	`)
	if err != nil {
		return nil, err
	}

	resp = &Response{}
	err = stmt.QueryRow(id).Scan(
		&resp.ID,
		&resp.ProjectID,
		&resp.Status,
		&resp.Length,
		&resp.Elapsed,
		&resp.Edited,
		&resp.Timestamp,
		&resp.Comment,
		&resp.Raw,
	)

	return resp, err
}
