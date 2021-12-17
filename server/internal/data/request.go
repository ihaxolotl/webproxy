package data

import (
	"database/sql"
	"time"
)

// Request represents an HTTP request and its metadata that has
// been intercepted by the proxy.
type Request struct {
	ID        string    // Unique ID of the request.
	ProjectID string    // Unique ID of the parent project.
	Method    string    // HTTP method of the request.
	Domain    string    // Domain name of the target host.
	IPAddr    string    // Internet address of the target host.
	Length    int64     // Length of the request in bytes.
	Edited    bool      // Flag for whether the request was modified or not.
	Timestamp time.Time // Time the request was made.
	Comment   string    // User-supplied comment on the request.
	Raw       string    // Raw request bytes.
}

// CreateRequestTable creates the "requests" table if it doesn't already exist.
func CreateRequestTable(db *sql.DB) (err error) {
	_, err = db.Exec(`
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

// InsertRequest inserts a new entry into the requests table and returns
// the last inserted rowid and an error.
func InsertRequest(db *sql.DB, req *Request) (rowid int64, err error) {
	var (
		stmt *sql.Stmt
		res  sql.Result
	)

	stmt, err = db.Prepare(`
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

func InsertAndGetRequest(db *sql.DB, r *Request) (req *Request, err error) {
	var (
		stmt  *sql.Stmt
		rowid int64
	)

	rowid, err = InsertRequest(db, r)
	if err != nil {
		return nil, err
	}

	stmt, err = db.Prepare(`
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

// GetRequestById queries and returns the entry from the requests table matching
// the supplied id parameter and optionally, an error.
func GetRequestById(db *sql.DB, id string) (req *Request, err error) {
	var stmt *sql.Stmt

	stmt, err = db.Prepare(`
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
