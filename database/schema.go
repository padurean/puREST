package database

import "database/sql"

var schema = `
CREATE TABLE IF NOT EXISTS person (
    first_name text,
    last_name text,
    email text
);

ALTER TABLE IF EXISTS ONLY person
	ADD COLUMN IF NOT EXISTS phone_number VARCHAR(255);

CREATE TABLE IF NOT EXISTS place (
    country text,
    city text NULL,
    telcode integer
)`

// Person ...
type Person struct {
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	Email       string
	PhoneNumber sql.NullString `db:"phone_number"`
}

// Place ...
type Place struct {
	Country string
	City    sql.NullString
	TelCode int
}
