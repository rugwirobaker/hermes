package sqlite

import "github.com/mattn/go-sqlite3"

func IsUniqueConstraintError(err error) bool {
	sqlite3Err, ok := err.(sqlite3.Error)
	if !ok {
		return false
	}
	return sqlite3Err.ExtendedCode == sqlite3.ErrConstraintUnique
}
