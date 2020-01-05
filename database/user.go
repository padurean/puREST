package database

import (
	"database/sql"
	"time"
)

// User ...
type User struct {
	ID        int64          `json:"id"`
	Username  string         `json:"username" validate:"required"`
	Password  string         `json:"password" validate:"required"`
	Email     string         `json:"email" validate:"required,email"`
	FirstName sql.NullString `json:"first_name" db:"first_name"`
	LastName  sql.NullString `json:"last_name" db:"last_name"`
	Created   time.Time      `json:"created"`
	Updated   time.Time      `json:"updated"`
}

var userSQLInsert string
var userSQLSelectByID string

func init() {
	userSQLInsert = `INSERT INTO ` + dbSchema + `.user (username, password, email, first_name, last_name)
		VALUES (:username, :password, :email, :first_name, :last_name) RETURNING id`
	userSQLSelectByID = `SELECT * FROM ` + dbSchema + `.user WHERE id=$1`
}

// Create ...
func (u *User) Create(db *DB) (*User, error) {
	var uu User
	if err := Create(db, userSQLInsert, userSQLSelectByID, u, &uu); err != nil {
		return nil, err
	}
	return &uu, nil
}

// Get ...
func (u *User) Get(db *DB) (*User, error) {
	var uu User
	if err := Get(db, userSQLSelectByID, u.ID, &uu); err != nil {
		return nil, err
	}
	return &uu, nil
}
