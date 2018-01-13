package models

import (
	// Bring in the SQLite3 functionality
	_ "github.com/mattn/go-sqlite3"
)

// User provides the user definition
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

// GetUsers returns all users
func (db *DB) GetUsers() ([]*User, error) {
	sqlQuery := `
	SELECT
		_rowid_,
		username,
		firstname,
		lastname,
		email
	FROM users
	`
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Firstname,
			&user.Lastname,
			&user.Email,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUser gets the user
func (db *DB) GetUser(id int) error {
	sqlQuery := `
	SELECT
		_rowid_,
		username,
		firstname,
		lastname,
		email
	FROM users
	WHERE _rowid_=?
	`
	rows, err := db.Query(sqlQuery, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	user := new(User)
	for rows.Next() {
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Firstname,
			&user.Lastname,
			&user.Email,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateUser creates a user
func (db *DB) CreateUser(user *User) error {
	sqlQuery := `
	INSERT INTO users (
		username,
		firstname,
		lastname,
		email
	) VALUES (?, ?, ?, ?)
	`

	res, err := db.Exec(
		sqlQuery,
		user.Username,
		user.Firstname,
		user.Lastname,
		user.Email,
	)
	if err != nil {
		return err
	}
	insertedID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(insertedID)
	return nil
}

// DeleteUser deletes a user
func (db *DB) DeleteUser(id int) error {
	sqlQuery := "DELETE FROM users WHERE _rowid_=?"

	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates a user
func (db *DB) UpdateUser(user *User) error {
	sqlQuery := `
	UPDATE users SET
		username=?,
		firstname=?,
		lastname=?,
		email=?
	WHERE _rowid_=?
	`

	_, err := db.Exec(
		sqlQuery,
		user.Username,
		user.Firstname,
		user.Lastname,
		user.Email,
		user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
