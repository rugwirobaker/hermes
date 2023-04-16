package hermes

import (
	"context"
	"time"

	"github.com/rugwirobaker/hermes/observ"
	"github.com/rugwirobaker/hermes/sqlite"
)

type Message struct {
	ID         int       `json:"id"`
	From       string    `json:"from"`
	ProviderID int       `json:"provider_id"`
	Recipient  string    `json:"phone"`
	Payload    string    `json:"payload"`
	Cost       float64   `json:"cost"`
	Status     Status    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdateAt   time.Time `json:"updated_at"`
}

type Store interface {
	// Insert a new message
	Insert(context.Context, *Message) (*Message, error)
	// MessageByID returns a message by serial id
	MessageBySerial(context.Context, string) (*Message, error)
	// MessageByPhone returns a message by phone
	MessageByPhone(context.Context, string) (*Message, error)
	// MessageByID returns a message by provider id
	MessageByID(context.Context, string) (*Message, error)
	// Update a message(set status to delivered/failed)
	Update(context.Context, *Message) (*Message, error)
}

func NewStore(db *sqlite.DB) Store {
	return &store{
		db: db,
	}
}

type store struct {
	db *sqlite.DB
}

func (s *store) Insert(ctx context.Context, u *Message) (*Message, error) {
	const op = "store.Insert"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	tx, err := s.db.BeginTx(ctx, sqlite.TxOptions(false))
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// scan

	res, err := tx.ExecContext(ctx, "INSERT INTO messages (provider_id, phone, payload, cost, status) VALUES (?, ?, ?, ?, ?)", u.ProviderID, u.Recipient, u.Payload, u.Cost, u.Status)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.ID = int(id)

	return u, tx.Commit()
}

func (s *store) MessageBySerial(ctx context.Context, id string) (*Message, error) {
	const op = "store.MessageBySerial"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	tx, err := s.db.BeginTx(ctx, sqlite.TxOptions(true))
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var out Message
	row := tx.QueryRowContext(ctx, "SELECT id, provider_id, phone, payload, cost, status, created_at, updated_at FROM messages WHERE id = ?", id)

	err = row.Scan(&out.ID,
		&out.ProviderID,
		&out.Recipient,
		&out.Payload,
		&out.Cost,
		&out.Status,
		&out.CreatedAt,
		&out.UpdateAt,
	)
	if err != nil {
		return nil, err
	}
	return &out, tx.Commit()
}

func (s *store) MessageByPhone(ctx context.Context, phone string) (*Message, error) {
	const op = "store.MessageByPhone"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	tx, err := s.db.BeginTx(ctx, sqlite.TxOptions(true))
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var out Message

	var row = tx.QueryRowContext(ctx, "SELECT id, phone, payload, cost, status, created_at, updated_at FROM messages WHERE phone = ?", phone)

	err = row.Scan(&out.ID,
		&out.Recipient,
		&out.Payload,
		&out.Cost,
		&out.Status,
		&out.CreatedAt,
		&out.UpdateAt,
	)
	if err != nil {
		return nil, err
	}
	return &out, tx.Commit()
}

func (s *store) MessageByID(ctx context.Context, id string) (*Message, error) {
	const op = "store.MessageByID"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	tx, err := s.db.BeginTx(ctx, sqlite.TxOptions(true))
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var out Message

	var row = tx.QueryRowContext(ctx, "SELECT id, provider_id, phone, payload, cost, status, created_at, updated_at FROM messages WHERE provider_id = ?", id)

	err = row.Scan(
		&out.ID,
		&out.ProviderID,
		&out.Recipient,
		&out.Payload,
		&out.Cost,
		&out.Status,
		&out.CreatedAt,
		&out.UpdateAt,
	)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *store) Update(ctx context.Context, u *Message) (*Message, error) {
	const op = "store.Update"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	tx, err := s.db.BeginTx(ctx, sqlite.TxOptions(false))
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE messages SET status = ? WHERE id = ?", u.Status, u.ID)
	if err != nil {
		return nil, err
	}

	return u, tx.Commit()
}
