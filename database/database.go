package database

import (
	"database/sql"
	"fmt"

	// initialize database (PosgreSQL) driver
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

// DB ...
type DB struct {
	*sqlx.DB
}

// MustConnect ...
func MustConnect(driver string, url string) *DB {
	return &DB{DB: sqlx.MustConnect(driver, url)}
}

// Create ...
func Create(db *DB, insert string, selectByID string, arg interface{}, dest interface{}) error {
	stmtInsert, err := db.PrepareNamed(insert)
	if err != nil {
		return fmt.Errorf("error preparing named db insert: %v", err)
	}
	var id int64
	if err := stmtInsert.Get(&id, arg); err != nil {
		return fmt.Errorf("error executing named db insert: %v", err)
	}
	stmtSelectByID, err := db.Preparex(selectByID)
	if err != nil {
		return fmt.Errorf("error preparing db select by ID: %v", err)
	}
	if err := stmtSelectByID.Get(dest, id); err != nil {
		return fmt.Errorf("error executing db select by ID: %v", err)
	}
	return nil
}

// Get ...
func Get(db *DB, selectByID string, id int64, dest interface{}) error {
	stmtSelectByID, err := db.Preparex(selectByID)
	if err != nil {
		return fmt.Errorf("error preparing db select by ID: %v", err)
	}
	if err := stmtSelectByID.Get(dest, id); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return fmt.Errorf("no record found for ID %d", id)
		default:
			return fmt.Errorf("error executing db select by ID: %v", err)
		}
	}
	return nil
}
