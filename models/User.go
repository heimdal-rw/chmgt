package models

import (
	"database/sql"

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
func GetUsers(db *sql.DB) ([]User, error) {
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

	users := []User{}
	for rows.Next() {
		var user User
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
func (user *User) GetUser(db *sql.DB) error {
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
	rows, err := db.Query(sqlQuery, user.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

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
func (user *User) CreateUser(db *sql.DB) error {
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
func (user *User) DeleteUser(db *sql.DB) error {
	sqlQuery := "DELETE FROM users WHERE _rowid_=?"

	_, err := db.Exec(sqlQuery, user.ID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates a user
func (user *User) UpdateUser(db *sql.DB) error {
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
	)
	if err != nil {
		return err
	}

	return nil
}
