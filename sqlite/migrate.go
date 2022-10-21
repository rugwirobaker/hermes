package sqlite

import (
	"database/sql"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
)

type Direction int

// Migration directions
const (
	// Migration apply
	Up Direction = 0
	// Migration Rollback
	Down Direction = 1
)

func Migrate(db *sql.DB, dir Direction, driver string) (int, error) {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				// id uses current unix timestamp
				Id: "1",
				Up: []string{
					createMessagesTable,
					createIndexOnProviderID,
				},
				DisableTransactionDown: true,
			},
			{
				Id: "2",
				Up: []string{
					createIdempotencyKeysTable,
					createUniqueIndexOnIdempotencyKey,
				},
				DisableTransactionDown: true,
			},
			{
				Id: "3",
				Up: []string{
					alterMessagesTableAddIdempotencyKeyReference,
					alterMessagesTableAddIdempotencyKeyReferenceIndex,
				},
				DisableTransactionDown: true,
			},
		},
	}

	n, err := migrate.Exec(db, driver, migrations, migrate.MigrationDirection(dir))
	if err != nil {
		return n, fmt.Errorf("could not apply migrations %w", err)
	}
	return n, nil
}

//flavor:sqlite3
const createMessagesTable = `CREATE TABLE IF NOT EXISTS messages (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	provider_id TEXT NOT NULL,
	phone TEXT NOT NULL,
	payload TEXT NOT NULL,
	cost INTEGER NOT NULL,
	status TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

//flavor:sqlite3
const createIndexOnProviderID = `CREATE INDEX IF NOT EXISTS idx_provider_id ON messages (provider_id);`

// flavor:sqlite3
const createIdempotencyKeysTable = `CREATE TABLE IF NOT EXISTS idempotency_keys (
	id INTEGER 		PRIMARY KEY AUTOINCREMENT,
	created_at 		DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	idempotency_key TEXT NOT NULL CHECK (length(idempotency_key) <= 100),
	last_run_at 	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    locked_at       DATETIME DEFAULT CURRENT_TIMESTAMP,

	/* parameters of the incoming request */
	request_method  TEXT 		NOT NULL CHECK (length(request_method) <= 10),
	request_path    TEXT 		NOT NULL CHECK (length(request_path) <= 100),
	request_params  BLOB 		NOT NULL,
	request_body    BLOB        NOT NULL,
	
	/* for finished requests, stored status code, and body */
	response_code   INT         NULL,
	response_body   BLOB        NULL,

	recovery_point  TEXT        NOT NULL CHECK (length(recovery_point) <= 50)
		
	/* may add a foreign key to the application(user) table later */
);
`
const createUniqueIndexOnIdempotencyKey = `CREATE UNIQUE INDEX IF NOT EXISTS idempotency_keys_key_idx ON idempotency_keys (idempotency_key);`

const alterMessagesTableAddIdempotencyKeyReference = `ALTER TABLE messages ADD COLUMN idempotency_key_id INTEGER REFERENCES idempotency_keys(id) ON DELETE SET NULL;`
const alterMessagesTableAddIdempotencyKeyReferenceIndex = `CREATE INDEX IF NOT EXISTS idx_idempotency_key_id ON messages (idempotency_key_id);`
