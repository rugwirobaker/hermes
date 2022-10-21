package hermes

import (
	"context"
	"time"

	"github.com/rugwirobaker/hermes/observ"
	"github.com/rugwirobaker/hermes/sqlite"
)

type RecoveryPoint string

const (
	RecoveryPointStart    RecoveryPoint = "started"
	RecoveryPointFinished RecoveryPoint = "finished"
)

func (r RecoveryPoint) String() string {
	return string(r)
}

type IdempotencyKey struct {
	ID        int       `json:"id"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LockedAt  time.Time `json:"locked_at"`

	RequestMethod string `json:"request_method"`
	RequestPath   string `json:"request_path"`
	RequestParams map[string]string
	RequestBody   []byte `json:"request_body"`
	RequestHeader map[string][]string

	ResponseCode int    `json:"response_code"`
	ResponseBody []byte `json:"response_body"`

	Recovery RecoveryPoint `json:"recovery"`
}

type IdempotencyKeyStore interface {
	// Insert idempotency key
	Create(context.Context, *IdempotencyKey) (*IdempotencyKey, error)
	// Key returns a key by key
	Key(context.Context, string) (*IdempotencyKey, error)
	// Update a key
	Update(context.Context, *IdempotencyKey) (*IdempotencyKey, error)
}

type idempotencyKeyStore struct {
	db *sqlite.DB
}

func NewIdempotencyKeyStore(db *sqlite.DB) IdempotencyKeyStore {
	return &idempotencyKeyStore{
		db: db,
	}
}

func (s *idempotencyKeyStore) Create(ctx context.Context, in *IdempotencyKey) (*IdempotencyKey, error) {
	const op = "idempotencyKeyStore.Insert"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	tx, err := s.db.BeginTx(ctx, sqlite.TxOptions(false))
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		`INSERT INTO idempotency_keys (
			key, 
			recovery
		) VALUES (?, ?)`, in.Key, in.Recovery)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	in.ID = int(id)
	in.CreatedAt = time.Now()
	in.UpdatedAt = time.Now()

	return in, nil
}

func (s *idempotencyKeyStore) Key(ctx context.Context, key string) (*IdempotencyKey, error) {
	const op = "idempotencyKeyStore.Key"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	tx, err := s.db.BeginTx(ctx, sqlite.TxOptions(false))
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// scan

	var u IdempotencyKey
	err = tx.QueryRowContext(ctx, "SELECT * FROM idempotency_keys WHERE key = ?", key).Scan(&u.ID, &u.Key, &u.CreatedAt, &u.UpdatedAt, &u.LockedAt, &u.RequestMethod, &u.RequestPath, &u.RequestParams, &u.RequestBody, &u.RequestHeader, &u.ResponseCode, &u.ResponseBody, &u.Recovery)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *idempotencyKeyStore) Update(ctx context.Context, in *IdempotencyKey) (*IdempotencyKey, error) {
	const op = "idempotencyKeyStore.Update"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	tx, err := s.db.BeginTx(ctx, sqlite.TxOptions(false))
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// scan

	res, err := tx.ExecContext(
		ctx, `
		UPDATE 
			idempotency_keys 
		SET 
			request_method = ?, 
			request_path = ?, 
			request_params = ?, 
			request_body = ?, 
			request_header = ?, 
			response_code = ?, 
			response_body = ?, 
			recovery_point = ? 
		WHERE 
			key = ?
	`,
		in.RequestMethod,
		in.RequestPath,
		in.RequestParams,
		in.RequestBody,
		in.RequestHeader,
		in.ResponseCode,
		in.ResponseBody,
		in.Recovery,
		in.Key,
	)

	if err != nil {
		return nil, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}

	in.UpdatedAt = time.Now()

	return in, nil
}

func (s *idempotencyKeyStore) Lock(ctx context.Context, key string) (*IdempotencyKey, error) {
	const op = "idempotencyKeyStore.Lock"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	tx, err := s.db.BeginTx(ctx, sqlite.TxOptions(false))
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// scan

	var u IdempotencyKey

	var row = tx.QueryRowContext(ctx, `
	SELECT 
		id,
		key,
		created_at,
		updated_at,
		locked_at,
		request_method,
		request_path,
		request_params,
		request_body,
		request_header,
		response_code,
		response_body,
		recovery_point
	FROM 
		idempotency_keys 
	WHERE key = ?`, key)

	err = row.Scan(&u.ID,
		&u.Key,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.LockedAt,
		&u.RequestMethod,
		&u.RequestPath,
		&u.RequestParams,
		&u.RequestBody,
		&u.RequestHeader,
		&u.ResponseCode,
		&u.ResponseBody,
		&u.Recovery,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
