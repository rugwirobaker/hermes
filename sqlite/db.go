package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rugwirobaker/hermes/observ"

	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

const Driver = "sqlite3"

type DB struct {
	db *sql.DB
}

func NewDB(path string, provider trace.TracerProvider) (*DB, error) {
	dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc&_journal_mode=WAL", path)

	db, err := otelsql.Open(Driver,
		dsn,
		otelsql.WithAttributes(semconv.DBSystemSqlite),
		otelsql.WithDBName(path),
		otelsql.WithTracerProvider(provider),
	)

	if err != nil {
		return nil, fmt.Errorf("could not open sqlite database %w", err)
	}

	// migrate database
	_, err = Migrate(db, Up)
	if err != nil {
		return nil, err
	}

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
