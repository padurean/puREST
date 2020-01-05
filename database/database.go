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

// Upsert ...
func Upsert(db *DB, sqlUpsert string, sqlSelectByID string, argUpsert interface{}, dest interface{}) error {
	stmtUpsert, err := db.PrepareNamed(sqlUpsert)
	if err != nil {
		return fmt.Errorf("error preparing named db upsert: %v", err)
	}
	var id int64
	if err := stmtUpsert.Get(&id, argUpsert); err != nil {
		return fmt.Errorf("error executing db upsert: %v", err)
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

// SelectOne ...
func SelectOne(db *DB, sqlSelect string, argSelect interface{}, dest interface{}) error {
	stmtSelect, err := db.Preparex(sqlSelect)
	if err != nil {
		return fmt.Errorf("error preparing db select one: %v", err)
	}
	return stmtSelect.Get(dest, argSelect)
}
