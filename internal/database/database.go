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

// LimitAndOffset ...
type LimitAndOffset struct {
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
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

// MarkAsDeleted ...
func MarkAsDeleted(db *DB, sqlMarkAsDeleted string, id int64) error {
	stmtMarkAsDeleted, err := db.Preparex(sqlMarkAsDeleted)
	if err != nil {
		return fmt.Errorf("error preparing db mark as deleted ID %d: %v", id, err)
	}
	result, err := stmtMarkAsDeleted.Exec(id)
	if err != nil {
		return fmt.Errorf("error executing db mark as deleted ID %d: %v", id, err)
	}
	nbMarkedAsDeleted, err := result.RowsAffected()
	if nbMarkedAsDeleted == 0 {
		return fmt.Errorf("error marking as deleted ID %d: nothing was marked as deleted", id)
	} else if nbMarkedAsDeleted > 1 {
		return fmt.Errorf("error marking as deleted ID %d: %d entities have been marked as deleted instead of just 1", id, nbMarkedAsDeleted)
	}
	return nil
}
