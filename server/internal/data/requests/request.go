package requests

import (
	"database/sql"
	"time"
)

// Request represents an HTTP request and its metadata that has
// been intercepted by the proxy.
type Request struct {
	ID        string    `json:"id"`        // Unique ID of the request.
	ProjectID string    `json:"projectId"` // Unique ID of the parent project.
	Method    string    `json:"method"`    // HTTP method of the request.
	Domain    string    `json:"domain"`    // Domain name of the target host.
	IPAddr    string    `json:"ipaddr"`    // Internet address of the target host.
	Length    int64     `json:"length"`    // Length of the request in bytes.
	Edited    bool      `json:"edited"`    // Flag for whether the request was modified or not.
	Timestamp time.Time `json:"timestamp"` // Time the request was made.
	Comment   string    `json:"comment"`   // User-supplied comment on the request.
	Raw       string    `json:"raw"`       // Raw request bytes.
}

type RequestsTable struct {
	db *sql.DB
}

func New(db *sql.DB) *RequestsTable {
	return &RequestsTable{db}
}

// Create creates the "requests" table if it doesn't already exist.
func (t RequestsTable) Create() (err error) {
	_, err = t.db.Exec(`
		CREATE TABLE IF NOT EXISTS requests (
			id TEXT PRIMARY KEY NOT NULL UNIQUE,
			projectid TEXT NOT NULL,
			method TEXT NOT NULL,
			domain TEXT NOT NULL,
			ipaddr TEXT NOT NULL,
			length INTEGER,
			edited BOOLEAN NOT NULL CHECK (edited IN (0, 1)),
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			comment TEXT,
			raw TEXT
		);
	`)

	return err
}

// InsertRequest inserts a new record into the requests table and returns
// the last inserted rowid or an error.
func (t RequestsTable) Insert(req *Request) (rowid int64, err error) {
	var (
		stmt *sql.Stmt
		res  sql.Result
	)

	stmt, err = t.db.Prepare(`
		INSERT INTO requests(
			id,
			projectid,
			method,
			domain,
			ipaddr,
			length,
			edited,
			timestamp,
			comment,
			raw
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		);
	`)
	if err != nil {
		return 0, err
	}

	res, err = stmt.Exec(
		req.ID,
		req.ProjectID,
		req.Method,
		req.Domain,
		req.IPAddr,
		req.Length,
		req.Edited,
		req.Timestamp,
		req.Comment,
		req.Raw,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// InsertAndFetch inserts a new record into the requests table and returns the
// full record or an error.
func (t RequestsTable) InsertAndFetch(r *Request) (req *Request, err error) {
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
			method,
			domain,
			ipaddr,
			length,
			edited,
			timestamp,
			comment,
			raw
		FROM
			requests
		WHERE
			rowid = ?
		LIMIT 0, 1;

	`)
	if err != nil {
		return nil, err
	}

	req = &Request{}
	err = stmt.QueryRow(rowid).Scan(
		&req.ID,
		&req.ProjectID,
		&req.Method,
		&req.Domain,
		&req.IPAddr,
		&req.Length,
		&req.Edited,
		&req.Timestamp,
		&req.Comment,
		&req.Raw,
	)

	return req, err
}

// FetchById queries and returns the record from the requests table matching
// the supplied id parameter and optionally, an error.
func (t RequestsTable) FetchById(id string) (req *Request, err error) {
	var stmt *sql.Stmt

	stmt, err = t.db.Prepare(`
		SELECT 
			id,
			projectid,
			method,
			domain,
			ipaddr,
			length,
			edited,
			timestamp,
			comment,
			raw
		FROM
			requests
		WHERE
			id = ?
		LIMIT 0, 1;

	`)
	if err != nil {
		return nil, err
	}

	req = &Request{}
	err = stmt.QueryRow(id).Scan(
		&req.ID,
		&req.ProjectID,
		&req.Method,
		&req.Domain,
		&req.IPAddr,
		&req.Length,
		&req.Edited,
		&req.Timestamp,
		&req.Comment,
		&req.Raw,
	)

	return req, err
}
