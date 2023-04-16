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
				Down: []string{"DROP TABLE messages;"},
			},
			{
				Id: "2",
				Up: []string{
					createAppsTable,
				},
			},
			{
				Id: "3",
				Up: []string{
					alterMessagesTableChangeCostToFloat,
				},
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

// hermes.App migration
const createAppsTable = `CREATE TABLE IF NOT EXISTS apps (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	token TEXT NOT NULL,
	sender TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	message_count INTEGER NOT NULL DEFAULT 0
);`

// change messages.cost to float64
const alterMessagesTableChangeCostToFloat = `CREATE TEMPORARY TABLE messages_temp (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	provider_id TEXT NOT NULL,
	phone TEXT NOT NULL,
	payload TEXT NOT NULL,
	cost REAL NOT NULL, -- Changed to REAL
	status TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO messages_temp (id, 
	provider_id, 
	phone, 
	payload, 
	cost, 
	status, 
	created_at, 
updated_at) SELECT id, provider_id, phone, payload, CAST(cost AS REAL), status, created_at, updated_at FROM messages;
DROP TABLE messages;
ALTER TABLE messages_temp RENAME TO messages;`
