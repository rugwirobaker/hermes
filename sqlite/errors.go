package sqlite

import (
	"database/sql"

	"github.com/mattn/go-sqlite3"
)

func IsUniqueConstraintError(err error) bool {
	sqlite3Err, ok := err.(sqlite3.Error)
	if !ok {
		return false
	}
	return sqlite3Err.ExtendedCode == sqlite3.ErrConstraintUnique
}

func IsForeignKeyConstraintError(err error) bool {
	sqlite3Err, ok := err.(sqlite3.Error)
	if !ok {
		return false
	}
	return sqlite3Err.ExtendedCode == sqlite3.ErrConstraintForeignKey
}

// invalid_text_representation
func IsInvalid(err error) bool {
	sqlite3Err, ok := err.(sqlite3.Error)
	if !ok {
		return false
	}
	return sqlite3Err.ExtendedCode == sqlite3.ErrConstraintCheck
}

// IsNoRowsError returns true if the error is a sql.ErrNoRows error.
func IsNoRowsError(err error) bool {
	return err == sql.ErrNoRows
}
