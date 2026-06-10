package ai

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// MemoryStore persists app context and conversation history in PostgreSQL.
type MemoryStore struct {
	db *sql.DB
}

// NewMemoryStore connects to the PostgreSQL instance on the NUCBox.
func NewMemoryStore(host string, port int, dbname, user, password, sslmode string) (*MemoryStore, error) {
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		host, port, dbname, user, password, sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &MemoryStore{db: db}, nil
}

// InitSchema creates the memory tables if they don't exist.
func (m *MemoryStore) InitSchema(ctx context.Context) error {
	_, err := m.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS kyanite_context (
			app       TEXT NOT NULL,
			session_id TEXT NOT NULL,
			key       TEXT NOT NULL,
			value     TEXT NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (app, session_id, key)
		);

		CREATE TABLE IF NOT EXISTS kyanite_conversations (
			id         SERIAL PRIMARY KEY,
			app        TEXT NOT NULL,
			session_id TEXT NOT NULL,
			role       TEXT NOT NULL,
			content    TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_conversations_app_session
			ON kyanite_conversations (app, session_id, created_at DESC);
	`)
	return err
}

// Close closes the database connection.
func (m *MemoryStore) Close() error {
	return m.db.Close()
}

// SaveContext stores a key-value pair for an app session.
func (m *MemoryStore) SaveContext(ctx context.Context, app, sessionID, key string, value any) error {
	valBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal value: %w", err)
	}

	_, err = m.db.ExecContext(ctx, `
		INSERT INTO kyanite_context (app, session_id, key, value, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (app, session_id, key)
		DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
	`, app, sessionID, key, string(valBytes))

	return err
}

// LoadContext retrieves a key's value for an app session.
func (m *MemoryStore) LoadContext(ctx context.Context, app, sessionID, key string, dest any) error {
	var val string
	err := m.db.QueryRowContext(ctx, `
		SELECT value FROM kyanite_context
		WHERE app = $1 AND session_id = $2 AND key = $3
	`, app, sessionID, key).Scan(&val)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// LoadAllContext retrieves all key-value pairs for an app session.
func (m *MemoryStore) LoadAllContext(ctx context.Context, app, sessionID string) (map[string]json.RawMessage, error) {
	rows, err := m.db.QueryContext(ctx, `
		SELECT key, value FROM kyanite_context
		WHERE app = $1 AND session_id = $2
	`, app, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]json.RawMessage)
	for rows.Next() {
		var k string
		var v json.RawMessage
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		result[k] = v
	}
	return result, rows.Err()
}

// SaveMessage appends a message to a conversation.
func (m *MemoryStore) SaveMessage(ctx context.Context, app, sessionID, role, content string) error {
	_, err := m.db.ExecContext(ctx, `
		INSERT INTO kyanite_conversations (app, session_id, role, content)
		VALUES ($1, $2, $3, $4)
	`, app, sessionID, role, content)
	return err
}

// LoadConversation loads recent messages from a conversation.
func (m *MemoryStore) LoadConversation(ctx context.Context, app, sessionID string, limit int) ([]Message, error) {
	rows, err := m.db.QueryContext(ctx, `
		SELECT role, content, created_at FROM kyanite_conversations
		WHERE app = $1 AND session_id = $2
		ORDER BY created_at DESC
		LIMIT $3
	`, app, sessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.Role, &m.Content, &m.Time); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}

	// Reverse to chronological order
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}

	return msgs, rows.Err()
}

// DeleteContext removes all context for a session.
func (m *MemoryStore) DeleteContext(ctx context.Context, app, sessionID string) error {
	_, err := m.db.ExecContext(ctx, `
		DELETE FROM kyanite_context WHERE app = $1 AND session_id = $2
	`, app, sessionID)
	return err
}

// IsAvailable checks if the database is reachable.
func (m *MemoryStore) IsAvailable(ctx context.Context) bool {
	return m.db.PingContext(ctx) == nil
}
