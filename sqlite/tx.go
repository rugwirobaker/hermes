package sqlite

import (
	"database/sql"
)

// Tx wraps the SQL Tx object to provide a timestamp at the start of the transaction.
type Tx struct {
	*sql.Tx
	db *DB
	// Now time.Time
}
