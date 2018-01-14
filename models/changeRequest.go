package models

import (
	// Bring in the SQLite3 functionality
	_ "github.com/mattn/go-sqlite3"
)

// ChangeRequest provides the change request definition
type ChangeRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	AuthorID    int    `json:"authorId"`
	RequesterID int    `json:"requesterId"`
	Description string `json:"description"`
	Reason      string `json:"reason"`
	Risk        string `json:"risk"`
	Steps       string `json:"steps"`
	Revert      string `json:"revert"`
}

// GetChangeRequests returns all change requests
func (db *DB) GetChangeRequests() ([]*ChangeRequest, error) {
	sqlQuery := `
	SELECT
		_rowid_,
		title,
		authorId,
		requesterId,
		description,
		reason,
		risk,
		steps,
		revert
	FROM changeRequest
	`
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	crs := make([]*ChangeRequest, 0)
	for rows.Next() {
		cr := new(ChangeRequest)
		err = rows.Scan(
			&cr.ID,
			&cr.Title,
			&cr.AuthorID,
			&cr.RequesterID,
			&cr.Description,
			&cr.Reason,
			&cr.Risk,
			&cr.Steps,
			&cr.Revert,
		)
		if err != nil {
			return nil, err
		}
		crs = append(crs, cr)
	}

	return crs, nil
}

// GetChangeRequest gets the change request
func (db *DB) GetChangeRequest(id int) (*ChangeRequest, error) {
	sqlQuery := `
	SELECT
		title,
		authorId,
		requesterId,
		description,
		reason,
		risk,
		steps,
		revert
	FROM changeRequest
	WHERE _rowid_=?
	`
	rows, err := db.Query(sqlQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cr := new(ChangeRequest)
	for rows.Next() {
		err = rows.Scan(
			&cr.Title,
			&cr.AuthorID,
			&cr.RequesterID,
			&cr.Description,
			&cr.Reason,
			&cr.Risk,
			&cr.Steps,
			&cr.Revert,
		)
		if err != nil {
			return nil, err
		}
	}
	return cr, nil
}

// CreateChangeRequest creates a change request
func (db *DB) CreateChangeRequest(cr *ChangeRequest) error {
	sqlQuery := `
	INSERT INTO changeRequest (
		title,
		authorId,
		requesterId,
		description,
		reason,
		risk,
		steps,
		revert
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := db.Exec(
		sqlQuery,
		cr.Title,
		cr.AuthorID,
		cr.RequesterID,
		cr.Description,
		cr.Reason,
		cr.Risk,
		cr.Steps,
		cr.Revert,
	)
	if err != nil {
		return err
	}
	insertedID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	cr.ID = int(insertedID)
	return nil
}

// DeleteChangeRequest deletes a change request
func (db *DB) DeleteChangeRequest(id int) error {
	sqlQuery := "DELETE FROM changeRequest WHERE _rowid_=?"

	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateChangeRequest updates a change request
func (db *DB) UpdateChangeRequest(cr *ChangeRequest) error {
	sqlQuery := `
	UPDATE changeRequest SET
		title=?,
		authorId=?,
		requesterId=?,
		description=?,
		reason=?,
		risk=?,
		steps=?,
		revert=?
	WHERE _rowid_=?
	`

	_, err := db.Exec(
		sqlQuery,
		cr.Title,
		cr.AuthorID,
		cr.RequesterID,
		cr.Description,
		cr.Reason,
		cr.Risk,
		cr.Steps,
		cr.Revert,
		cr.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
