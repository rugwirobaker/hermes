package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/rugwirobaker/hermes/observ"
)

func DeleteOldRecords(ctx context.Context, db *DB, ret time.Duration) (int64, error) {
	const op = "CleanUpIdempotencyKeys"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	tx, err := db.BeginTx(ctx, TxOptions(false))
	if err != nil {
		return 0, fmt.Errorf("error deleting old records: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, deleteIdempotencyKeys, ret.String())
	if err != nil {
		return 0, fmt.Errorf("error deleting old records: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error deleting old records: %w", err)
	}
	return rows, tx.Commit()
}

// deletes all idempotency keys older than 2 hours
var deleteIdempotencyKeys = `DELETE FROM idempotency_keys WHERE created_at < DATETIME('now', ?);`
