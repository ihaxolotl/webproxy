package history

import (
	"database/sql"
	"time"
)

type HistoryEntry struct {
	Index      int64     `json:"idx"`        // Request order index
	Method     string    `json:"method"`     // HTTP request method
	Status     int16     `json:"status"`     // HTTP response status code
	Target     string    `json:"target"`     // Target domain
	URL        string    `json:"url"`        // URL of the requested resource
	IPAddr     string    `json:"ipaddr"`     // Internet address of the target
	Length     int64     `json:"length"`     // Length of the response in bytes
	Timestamp  time.Time `json:"timestamp"`  // Timestamp of when the request was made
	Edited     bool      `json:"edited"`     // Flag for whether the request was modified
	Comment    string    `json:"comment"`    // Comment for the request
	RequestId  string    `json:"requestId"`  // Unique ID of the request
	ResponseId string    `json:"responseId"` // Unique ID of the response
}

type HistoryView struct {
	db *sql.DB
}

func New(db *sql.DB) *HistoryView {
	return &HistoryView{db}
}

func (v HistoryView) Create() (err error) {
	return nil
}

func (v HistoryView) Fetch(projectId string) (history []HistoryEntry, err error) {
	var (
		stmt *sql.Stmt
		rows *sql.Rows
	)

	stmt, err = v.db.Prepare(`
		SELECT
			ROW_NUMBER () OVER ( ORDER BY req.timestamp ) idx,
			req.method as method,
			res.status as status,
			req.domain as target,
			req.url as url,
			req.ipaddr as ipaddr,
			res.length as length,
			req.timestamp as timestamp,
			req.edited as edited,
			req.comment as comment,
			req.id as requestid,
			res.id as responseid
		FROM
			requests req
		INNER JOIN
			responses res
		ON
			req.id = res.requestid
		INNER JOIN
			projects proj
		ON
			proj.id = res.projectid
		WHERE
			proj.id = ?;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	history = make([]HistoryEntry, 0)

	rows, err = stmt.Query(projectId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var h HistoryEntry

		if err = rows.Scan(
			&h.Index,
			&h.Method,
			&h.Status,
			&h.Target,
			&h.URL,
			&h.IPAddr,
			&h.Length,
			&h.Timestamp,
			&h.Edited,
			&h.Comment,
			&h.RequestId,
			&h.ResponseId,
		); err != nil {
			return nil, err
		}

		history = append(history, h)
	}

	return history, err

}
