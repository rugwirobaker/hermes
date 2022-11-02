package sqlite

import (
	"database/sql"
	"log"

	"github.com/avast/retry-go"
	"github.com/mattn/go-sqlite3"
	"github.com/rugwirobaker/hermes/observ"
	"go.opentelemetry.io/otel/trace"
)

// Tx wraps the SQL Tx object to provide a timestamp at the start of the transaction.
type Tx struct {
	*sql.Tx
	db   *DB
	span trace.Span
	// Now time.Time
}

// Commit commits the transaction.
func (tx *Tx) Commit() (err error) {
	attempts, err := Retry(func() error {
		return tx.Tx.Commit()
	})
	if tx.span.IsRecording() {
		tx.span.SetAttributes(
			observ.Int64("attempts", int64(attempts)),
		)
	}
	return
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() error {
	return tx.Tx.Rollback()
}

func Retry(fn func() error) (attempts int, err error) {
	err = retry.Do(
		fn,
		retry.RetryIf(func(err error) bool {
			if e, ok := err.(sqlite3.Error); ok {
				return e.Code == sqlite3.ErrBusy || e.Code == sqlite3.ErrLocked
			}
			return false
		}),
		retry.Attempts(3),
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(100),
		retry.OnRetry(func(n uint, err error) {
			// record the number of retries
			attempts = int(n)
			log.Println("retrying commit", "attempt", n, "error", err)
		}),
	)
	return
}
