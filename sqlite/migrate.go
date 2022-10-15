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

func Migrate(db *sql.DB, dir Direction) (int, error) {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				// id uses current unix timestamp
				Id:   "1",
				Up:   []string{createMessagesTable},
				Down: []string{"DROP TABLE messages;"},
			},
		},
	}

	n, err := migrate.Exec(db, Driver, migrations, migrate.MigrationDirection(dir))
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
	status TEXT NOT NULL
);`
