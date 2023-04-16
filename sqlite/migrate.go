package sqlite

import (
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

func (db *DB) Migrate(dir Direction) (int, error) {
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
			{
				Id: "4",
				Up: []string{
					recreateMessagesTable,
				},
			},
			{
				Id: "5",
				Up: []string{
					alterMessagesTableChangeProviderIDToInteger,
				},
			},
		},
	}

	n, err := migrate.Exec(db.db, driver, migrations, migrate.MigrationDirection(dir))
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

// recreate messages table since I've lost my data anyway
const recreateMessagesTable = `CREATE TABLE IF NOT EXISTS messages (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	provider_id TEXT NOT NULL,
	phone TEXT NOT NULL,
	payload TEXT NOT NULL,
	cost REAL NOT NULL,
	status TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

// change the data type of messages.provider_id to INTEGER
const alterMessagesTableChangeProviderIDToInteger = `CREATE TABLE messages_temp (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	provider_id INTEGER NOT NULL, -- Changed to INTEGER
	phone TEXT NOT NULL,
	payload TEXT NOT NULL,
	cost REAL NOT NULL,
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
	updated_at) SELECT id, CAST(provider_id AS INTEGER), phone, payload, cost, status, created_at, updated_at FROM messages;
DROP TABLE messages;
ALTER TABLE messages_temp RENAME TO messages;`
