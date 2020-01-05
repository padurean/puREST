package controllers

import (
	"database/sql"
	"encoding/json"
)

// ContextKey ...
type ContextKey string

// ContextKey ...
const (
	ContextKeyUser ContextKey = "user"
	// ...
)

// NullString is a wrapper around sql.NullString
type NullString sql.NullString

// MarshalJSON method is called by json.Marshal,
// whenever it is of type NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON method is called by json.Unmarshal,
// whenever it is of type NullString
func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		ns.String = ""
		return nil
	}
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	ns.Valid = true
	ns.String = v
	return nil
}
