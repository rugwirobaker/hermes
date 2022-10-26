package hermes

import (
	"context"
	"database/sql"

	"github.com/rugwirobaker/hermes/observ"
	"github.com/rugwirobaker/hermes/sqlite"
)

// hermes is a service that exposes an API for sending SMS messages to other internal apps(services)

type App struct {
	// Name of the app
	Name string `json:"name"`
	// ID of the app
	ID string `json:"id"`
	// APIKey of the app
	APIKey string `json:"token"`
	// Sender is the sender of the message
	Sender string `json:"sender"`
	// CreatedAt is the time the app was created
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the time the app was updated
	UpdatedAt string `json:"updated_at"`
	// MessageCount is the number of messages sent by the app
	MessageCount int64 `json:"message_count"`
}

type AppStore interface {
	// Register a new app
	Register(context.Context, *App) error
	// Get an app by id
	Get(context.Context, string) (*App, error)
	// FindByToken finds an app by token
	FindByToken(context.Context, string) (*App, error)
	// Update an app
	Update(context.Context, *App) error
	// Delete an app
	Delete(context.Context, string) error
	// List all apps
	List(context.Context) ([]*App, error)
}

func NewAppStore(db *sqlite.DB) AppStore {
	return &appStore{db: db}
}

type appStore struct {
	db *sqlite.DB
}

func (s *appStore) Register(ctx context.Context, app *App) error {
	const op = "appStore.Register"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	opts := sqlite.TxOptions(false)

	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = tx.QueryRowContext(
		ctx,
		`INSERT INTO apps (
			name, 
			token, 
			sender
		) VALUES (?, ?, ?) RETURNING id, created_at, updated_at`,
		app.Name, app.APIKey, app.Sender,
	).Scan(&app.ID, &app.CreatedAt, &app.UpdatedAt); err != nil {
		if sqlite.IsUniqueConstraintError(err) {
			return ErrAlreadyExists
		}
		return err
	}

	return tx.Commit()
}

func (s *appStore) Get(ctx context.Context, id string) (*App, error) {
	const op = "appStore.Get"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	var app App

	opts := sqlite.TxOptions(true)

	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var row = tx.QueryRowContext(ctx, `
		SELECT 
			id, name, token, sender, created_at, updated_at, message_count
		FROM apps WHERE id = ?
	`, id)

	err = row.Scan(&app.ID, &app.Name, &app.APIKey, &app.Sender, &app.CreatedAt, &app.UpdatedAt, &app.MessageCount)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &app, tx.Commit()
}

func (s *appStore) Update(ctx context.Context, app *App) error {
	const op = "appStore.Update"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	opts := sqlite.TxOptions(false)

	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE apps
		SET name = ?, token = ?, sender = ?, created_at = ?, updated_at = ?, message_count = ?
		WHERE id = ?
	`, app.Name, app.APIKey, app.Sender, app.CreatedAt, app.UpdatedAt, app.MessageCount, app.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *appStore) Delete(ctx context.Context, id string) error {
	const op = "appStore.Delete"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	opts := sqlite.TxOptions(false)

	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		DELETE FROM apps
		WHERE id = ?
	`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *appStore) List(ctx context.Context) ([]*App, error) {
	const op = "appStore.List"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	var apps []*App

	opts := sqlite.TxOptions(true)

	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
		SELECT 
			id, 
			name, 
			token, 
			sender, 
			created_at, 
			updated_at, 
			message_count
		FROM apps
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var app App

		err = rows.Scan(&app.ID, &app.Name, &app.APIKey, &app.Sender, &app.CreatedAt, &app.UpdatedAt, &app.MessageCount)
		if err != nil {
			return nil, err
		}
		apps = append(apps, &app)
	}

	return apps, tx.Commit()
}

func (s *appStore) FindByToken(ctx context.Context, token string) (*App, error) {
	const op = "appStore.FindByToken"

	ctx, span := observ.StartSpan(ctx, op)
	defer span.End()

	var app App

	opts := sqlite.TxOptions(true)

	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var row = tx.QueryRowContext(ctx, `
		SELECT 
			id, name, token, sender, created_at, updated_at, message_count
		FROM apps WHERE token = ?
	`, token)

	err = row.Scan(&app.ID, &app.Name, &app.APIKey, &app.Sender, &app.CreatedAt, &app.UpdatedAt, &app.MessageCount)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &app, tx.Commit()
}
