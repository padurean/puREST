package database

import (
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
func Create(
	db *DB,
	sqlInsert string,
	sqlSelectByID string,
	argInsert interface{},
	dest interface{}) error {
	stmtInsert, err := db.PrepareNamed(sqlInsert)
	if err != nil {
		return fmt.Errorf("error preparing named db insert: %v", err)
	}
	var id int64
	if err := stmtInsert.Get(&id, argInsert); err != nil {
		return fmt.Errorf("error executing named db insert: %v", err)
	}
	stmtSelectByID, err := db.Preparex(sqlSelectByID)
	if err != nil {
		return fmt.Errorf("error preparing db select by ID: %v", err)
	}
	if err := stmtSelectByID.Get(dest, id); err != nil {
		return fmt.Errorf("error executing db select by ID: %v", err)
	}
	return nil
}

// Get ...
func Get(db *DB, sqlSelect string, argSelect interface{}, dest interface{}) error {
	stmtSelect, err := db.Preparex(sqlSelect)
	if err != nil {
		return fmt.Errorf("error preparing db select by ID: %v", err)
	}
	return stmtSelect.Get(dest, argSelect)
}
