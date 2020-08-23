package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/padurean/purest/internal/auth"
)

// User ...
type User struct {
	ID        int64          `json:"id"`
	Username  string         `json:"username" validate:"required,alphanum"`
	Password  string         `json:"password" validate:"required,password"`
	Email     string         `json:"email" validate:"required,email"`
	FirstName sql.NullString `json:"first_name" db:"first_name"`
	LastName  sql.NullString `json:"last_name" db:"last_name"`
	Role      auth.Role      `json:"role" validate:"required,role"`
	Created   time.Time      `json:"created"`
	Updated   time.Time      `json:"updated"`
	Deleted   sql.NullTime   `json:"deleted,omitempty"`
}

var userSQLInsert string
var userSQLUpdate string
var userSQLSelectByID string
var userSQLSelectByUsername string
var userSQLSelectByEmail string
var userSQLSelectList string
var userSQLMarkAsDeleted string

func init() {
	userSQLInsert = `INSERT INTO ` + dbSchema + `.user (username, password, email, first_name, last_name, role)
		VALUES (:username, :password, :email, :first_name, :last_name, :role) RETURNING id`
	userSQLUpdate = `UPDATE ` + dbSchema + `.user
		SET username=:username, password=:password, email=:email, first_name=:first_name, last_name=:last_name, role=:role, updated=CURRENT_TIMESTAMP
		WHERE id=:id RETURNING id`
	userSQLSelectByID = `SELECT * FROM ` + dbSchema + `.user WHERE id=$1`
	userSQLSelectByUsername = `SELECT * FROM ` + dbSchema + `.user WHERE username=$1`
	userSQLSelectByEmail = `SELECT * FROM ` + dbSchema + `.user WHERE email=$1`
	userSQLSelectList = `SELECT * FROM ` + dbSchema + `.user WHERE deleted IS NULL LIMIT :limit OFFSET :offset`
	userSQLMarkAsDeleted = `UPDATE ` + dbSchema + `.user SET deleted=CURRENT_TIMESTAMP WHERE id=$1 RETURNING id`
}

func (u *User) validateNoDuplicate(db *DB) error {
	usernameExists := true
	uWithSameUsername, err := u.GetByUsername(db)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			usernameExists = false
		default:
			return fmt.Errorf("error finding if an user with username %s already exists: %v", u.Username, err)
		}
	}
	if usernameExists && uWithSameUsername.ID != u.ID {
		return &ErrDuplicateRow{ColName: "username", ColValue: u.Username}
	}

	emailExists := true
	uWithSameEmail, err := u.GetByEmail(db)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			emailExists = false
		default:
			return fmt.Errorf("error finding if an user with email %s already exists: %v", u.Email, err)
		}
	}
	if emailExists && uWithSameEmail.ID != u.ID {
		return &ErrDuplicateRow{ColName: "email", ColValue: u.Email}
	}

	return nil
}

// Create ...
func (u *User) Create(db *DB) (*User, error) {
	if err := u.validateNoDuplicate(db); err != nil {
		return nil, err
	}

	var uu User
	if err := Upsert(db, userSQLInsert, userSQLSelectByID, u, &uu); err != nil {
		return nil, err
	}
	return &uu, nil
}

// Update ...
func (u *User) Update(db *DB) (*User, error) {
	if err := u.validateNoDuplicate(db); err != nil {
		return nil, err
	}

	var uu User
	if err := Upsert(db, userSQLUpdate, userSQLSelectByID, u, &uu); err != nil {
		return nil, err
	}
	return &uu, nil
}

// GetByID ...
func (u *User) GetByID(db *DB) (*User, error) {
	var uu User
	if err := SelectOne(db, userSQLSelectByID, u.ID, &uu); err != nil {
		return nil, err
	}
	return &uu, nil
}

// GetByUsername ...
func (u *User) GetByUsername(db *DB) (*User, error) {
	var uu User
	if err := SelectOne(db, userSQLSelectByUsername, u.Username, &uu); err != nil {
		return nil, err
	}
	return &uu, nil
}

// GetByEmail ...
func (u *User) GetByEmail(db *DB) (*User, error) {
	var uu User
	if err := SelectOne(db, userSQLSelectByEmail, u.Email, &uu); err != nil {
		return nil, err
	}
	return &uu, nil
}

// List ...
func (u *User) List(db *DB, limit int, offset int) ([]*User, error) {
	users := []*User{}
	usersRows, err := db.NamedQuery(userSQLSelectList, LimitAndOffset{Limit: limit, Offset: offset})
	if err != nil {
		return users, fmt.Errorf("error preparing named db select: %v", err)
	}
	defer usersRows.Close()
	for usersRows.Next() {
		u := User{}
		err := usersRows.StructScan(&u)
		if err != nil {
			return users, fmt.Errorf("error scanning user row to struct: %v", err)
		}
		users = append(users, &u)
	}
	err = usersRows.Close()
	if err != nil {
		return users, fmt.Errorf("error closing rows: %v", err)
	}

	return users, nil
}

// Delete ...
func (u *User) Delete(db *DB) error {
	return MarkAsDeleted(db, userSQLMarkAsDeleted, u.ID)
}
