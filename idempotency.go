package hermes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rugwirobaker/hermes/observ"
	"github.com/rugwirobaker/hermes/sqlite"
)

type IdempotencyRecord struct {
	Key     string
	Code    int
	Headers map[string][]string
	Body    []byte
	Path    string
}

// IdempotencyKeyStore is an interface for storing idempotency keys
type IdempotencyKeyStore interface {
	// Get returns an entry from the store
	Get(ctx context.Context, key string) (*IdempotencyRecord, error)
	// Set sets an entry in the store
	Set(ctx context.Context, entry *IdempotencyRecord) error
}

// idempotencyKeyStore is an sqlite backed store for idempotency keys
type idempotencyKeyStore struct {
	db *sqlite.DB
}

// NewIdempotencyKeyStore returns a new instance of idempotencyKeyStore
func NewIdempotencyKeyStore(db *sqlite.DB) IdempotencyKeyStore {
	return &idempotencyKeyStore{
		db: db,
	}
}

// Get returns an entry from the store
func (s *idempotencyKeyStore) Get(ctx context.Context, key string) (*IdempotencyRecord, error) {
	const op = "idempotencyKeyStore.Get"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	tx, err := s.db.BeginTx(ctx, sqlite.TxOptions(true))
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	defer tx.Rollback()

	var out = new(IdempotencyRecord)

	var headers string
	// remember the headers are stored as json we have to convert them back to map[string][]string
	err = tx.QueryRowContext(ctx, selectIdempotencyKey, key).Scan(&out.Key, &out.Code, &headers, &out.Body, &out.Path)
	if err != nil {
		if sqlite.IsNoRowsError(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(headers), &out.Headers)
	if err != nil {
		return nil, fmt.Errorf("insert: failed to unmarshal headers: %w", err)
	}
	return out, nil
}

// Set sets an entry in the store
func (s *idempotencyKeyStore) Set(ctx context.Context, entry *IdempotencyRecord) error {
	const op = "idempotencyKeyStore.Set"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	tx, err := s.db.BeginTx(ctx, sqlite.TxOptions(false))
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	defer tx.Rollback()

	headers, err := json.Marshal(entry.Headers)
	if err != nil {
		return fmt.Errorf("insert: failed to marshal headers: %w", err)
	}

	_, err = tx.ExecContext(ctx, insertIdempotencyKey, entry.Key, entry.Code, headers, entry.Body, entry.Path)
	if err != nil {
		if sqlite.IsUniqueConstraintError(err) || sqlite.IsForeignKeyConstraintError(err) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("insert: %w", err)
	}
	return tx.Commit()
}

var insertIdempotencyKey = `INSERT INTO idempotency_keys (key, code, headers, body, path) VALUES (?, ?, ?, ?, ?)`

var selectIdempotencyKey = `SELECT key, code, headers, body, path FROM idempotency_keys WHERE key = ?`
