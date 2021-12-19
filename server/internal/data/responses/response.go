package responses

import (
	"database/sql"
	"time"
)

// Response represents an HTTP response and its metadata that has
// been intercepted by the proxy.
type Response struct {
	ID        string    `json:"id"`        // Unique ID of the response.
	ProjectID string    `json:"projectId"` // Unique ID of the parent project.
	Status    int16     `json:"status"`    // HTTP status code of the response.
	Length    int64     `json:"length"`    // Length of the response in bytes.
	Elapsed   int64     `json:"elapsed"`   // Time elapsed since request was sent until response.
	Edited    bool      `json:"edited"`    // Flag for whether the response was modified or not.
	Timestamp time.Time `json:"timestamp"` // Time the response was received.
	Mimetype  string    `json:"mimetype"`  // Mime-type of the response body data.
	Comment   string    `json:"comment"`   // User-supplied comment on the response.
	Raw       string    `json:"raw"`       // Raw response bytes.
}

type ResponseTable struct {
	db *sql.DB
}

func New(db *sql.DB) *ResponseTable {
	return &ResponseTable{db}
}

// Create creates the "responses" table if it doesn't already exist.
func (t ResponseTable) Create() (err error) {
	_, err = t.db.Exec(`
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

// InsertResponse inserts a new record into the response table and returns
// the last inserted rowid or an error.
func (t ResponseTable) Insert(resp *Response) (rowid int64, err error) {
	var (
		stmt *sql.Stmt
		res  sql.Result
	)

	stmt, err = t.db.Prepare(`
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

// InsertAndFetch inserts a new record into the responses table and returns the
// full record or an error.
func (t ResponseTable) InsertAndFetch(r *Response) (resp *Response, err error) {
	var (
		stmt  *sql.Stmt
		rowid int64
	)

	rowid, err = t.Insert(r)
	if err != nil {
		return nil, err
	}

	stmt, err = t.db.Prepare(`
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

// FetchById queries and returns the record from the responses table matching
// the supplied id parameter and optionally, an error.
func (t ResponseTable) FetchById(id string) (resp *Response, err error) {
	var stmt *sql.Stmt

	stmt, err = t.db.Prepare(`
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
