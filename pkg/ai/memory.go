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

		CREATE TABLE IF NOT EXISTS kyanite_sessions (
			app        TEXT NOT NULL,
			session_id TEXT NOT NULL,
			title      TEXT,
			state      JSONB NOT NULL DEFAULT '{}',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(app, session_id)
		);

		CREATE INDEX IF NOT EXISTS idx_sessions_app
			ON kyanite_sessions (app, updated_at DESC);

		CREATE TABLE IF NOT EXISTS kyanite_cross_app_context (
			id              SERIAL PRIMARY KEY,
			source_app      TEXT NOT NULL,
			target_app      TEXT NOT NULL,
			context_type    TEXT NOT NULL,
			summary         TEXT NOT NULL,
			relevance_score REAL DEFAULT 0.5,
			created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_cross_app_target
			ON kyanite_cross_app_context (target_app, relevance_score DESC);
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

// ── Session Management ────────────────────────────────────────────────

// Session represents a saved app session for "resume where you left off".
type Session struct {
	App       string    `json:"app"`
	SessionID string    `json:"session_id"`
	Title     string    `json:"title"`
	State     string    `json:"state"` // JSON-encoded app state
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SaveSession persists an app session snapshot. Called on app exit.
func (m *MemoryStore) SaveSession(ctx context.Context, app, sessionID, title string, state any) error {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal session state: %w", err)
	}
	_, err = m.db.ExecContext(ctx, `
		INSERT INTO kyanite_sessions (app, session_id, title, state, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (app, session_id)
		DO UPDATE SET title = EXCLUDED.title, state = EXCLUDED.state, updated_at = NOW()
	`, app, sessionID, title, string(stateBytes))
	return err
}

// LoadSession restores a saved app session.
func (m *MemoryStore) LoadSession(ctx context.Context, app, sessionID string) (*Session, error) {
	var s Session
	var stateBytes string
	err := m.db.QueryRowContext(ctx, `
		SELECT app, session_id, COALESCE(title, ''), state, created_at, updated_at
		FROM kyanite_sessions
		WHERE app = $1 AND session_id = $2
	`, app, sessionID).Scan(&s.App, &s.SessionID, &s.Title, &stateBytes, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	s.State = stateBytes
	return &s, nil
}

// GetRecentSessions returns the N most recent sessions for an app.
func (m *MemoryStore) GetRecentSessions(ctx context.Context, app string, limit int) ([]Session, error) {
	rows, err := m.db.QueryContext(ctx, `
		SELECT app, session_id, COALESCE(title, ''), state, created_at, updated_at
		FROM kyanite_sessions
		WHERE app = $1
		ORDER BY updated_at DESC
		LIMIT $2
	`, app, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.App, &s.SessionID, &s.Title, &s.State, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

// GetAllRecentSessions returns the N most recent sessions across all apps.
func (m *MemoryStore) GetAllRecentSessions(ctx context.Context, limit int) ([]Session, error) {
	rows, err := m.db.QueryContext(ctx, `
		SELECT app, session_id, COALESCE(title, ''), state, created_at, updated_at
		FROM kyanite_sessions
		ORDER BY updated_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.App, &s.SessionID, &s.Title, &s.State, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

// CrossAppContext represents a piece of context from one app relevant to another.
type CrossAppContext struct {
	SourceApp      string    `json:"source_app"`
	TargetApp      string    `json:"target_app"`
	ContextType    string    `json:"context_type"`
	Summary        string    `json:"summary"`
	RelevanceScore float32   `json:"relevance_score"`
	CreatedAt      time.Time `json:"created_at"`
}

// SaveCrossAppContext stores a context link between apps.
func (m *MemoryStore) SaveCrossAppContext(ctx context.Context, c CrossAppContext) error {
	_, err := m.db.ExecContext(ctx, `
		INSERT INTO kyanite_cross_app_context (source_app, target_app, context_type, summary, relevance_score, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`, c.SourceApp, c.TargetApp, c.ContextType, c.Summary, c.RelevanceScore)
	return err
}

// GetCrossAppContext retrieves context from other apps relevant to the target app.
func (m *MemoryStore) GetCrossAppContext(ctx context.Context, targetApp string, limit int) ([]CrossAppContext, error) {
	rows, err := m.db.QueryContext(ctx, `
		SELECT source_app, target_app, context_type, summary, relevance_score, created_at
		FROM kyanite_cross_app_context
		WHERE target_app = $1
		ORDER BY relevance_score DESC, created_at DESC
		LIMIT $2
	`, targetApp, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CrossAppContext
	for rows.Next() {
		var c CrossAppContext
		if err := rows.Scan(&c.SourceApp, &c.TargetApp, &c.ContextType, &c.Summary, &c.RelevanceScore, &c.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	return results, rows.Err()
}
