package hermes

import (
	"context"

	"github.com/rugwirobaker/hermes/observ"
	"github.com/rugwirobaker/hermes/sqlite"
)

type Message struct {
	ID         int    `json:"id"`
	ProviderID string `json:"provider_id"`
	Recipient  string `json:"phone"`
	Payload    string `json:"payload"`
	Cost       int64  `json:"cost"`
	Status     Status `json:"status"`
}

type Store interface {
	// Insert a new message
	Insert(context.Context, *Message) (*Message, error)
	// MessageByID returns a message by id
	MessageByID(context.Context, string) (*Message, error)
	// MessageByPhone returns a message by phone
	MessageByPhone(context.Context, string) (*Message, error)
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
	row := tx.QueryRowContext(ctx, "SELECT id, provider_id, phone, payload, cost, status FROM messages WHERE id = ?", id)

	if err = row.Scan(&out.ID, &out.ProviderID, &out.Recipient, &out.Payload, &out.Cost, &out.Status); err != nil {
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

	err = tx.QueryRowContext(ctx, "SELECT id, phone, payload, cost, status FROM messages WHERE phone = ?", phone).Scan(&out.ID, &out.Recipient, &out.Payload, &out.Cost, &out.Status)
	if err != nil {
		return nil, err
	}
	return &out, tx.Commit()
}
