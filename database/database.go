package database

import (
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

// Migrate ...
func Migrate(db *DB) {
	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	db.MustExec(schema)
}
