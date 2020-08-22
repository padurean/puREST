package database

import "fmt"

// ErrDuplicateRow ...
type ErrDuplicateRow struct {
	ColName  string
	ColValue string
}

func (err *ErrDuplicateRow) Error() string {
	return fmt.Sprintf("%s '%s' already exists", err.ColName, err.ColValue)
}
