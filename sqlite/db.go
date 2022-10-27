package sqlite

import (
	"context"
	"database/sql"

	// _ "github.com/lib/pq"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rugwirobaker/hermes/observ"
	"github.com/rugwirobaker/hermes/tracing"
	"go.opentelemetry.io/otel/trace"
)

type DB struct {
	db *sql.DB
}

func NewDB(dsn, driver, server string, provider trace.TracerProvider) (*DB, error) {

	db, err := sql.Open(driver, dsn)

	if err != nil {
		return nil, err
	}

	// migrate database
	_, err = Migrate(db, Up, driver)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	// Register the otelsql wrapper for the provided database driver.
	driver, err = tracing.DBTraceDriver(provider, driver, dsn, server)

	if err != nil {
		return nil, err
	}

	// reopendb with otelsql wrapper
	db, err = sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// BeginTx starts a transaction and returns a wrapper Tx type. This type
// provides a reference to the database and a fixed timestamp at the start of
// the transaction. The timestamp allows us to mock time during tests as well.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	const op = "sqlite.BeginTx"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Return wrapper Tx that includes the transaction start time.
	return &Tx{
		Tx: tx,
		db: db,
	}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	// Close database.
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

func TxOptions(readonly bool) *sql.TxOptions {
	return &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  readonly,
	}
}
