package db

// Schema contains all database table creation statements
const Schema = `
-- Songs table (metadata index)
CREATE TABLE IF NOT EXISTS songs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filepath TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    artist TEXT,
    key TEXT,
    tempo INTEGER,
    time_signature TEXT,
    structure TEXT,
    tags TEXT, -- JSON array
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    word_count INTEGER DEFAULT 0,
    verse_count INTEGER DEFAULT 0,
    chorus_count INTEGER DEFAULT 0,
    quality_score REAL DEFAULT 0.0
);

-- Version history (auto-save snapshots)
CREATE TABLE IF NOT EXISTS versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    song_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    is_milestone BOOLEAN DEFAULT FALSE,
    milestone_name TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
);

-- Writing statistics
CREATE TABLE IF NOT EXISTS writing_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE NOT NULL UNIQUE,
    words_written INTEGER DEFAULT 0,
    songs_created INTEGER DEFAULT 0,
    songs_edited INTEGER DEFAULT 0,
    ai_requests INTEGER DEFAULT 0,
    time_spent_minutes INTEGER DEFAULT 0
);

-- Knowledge Base vectors (RAG)
CREATE TABLE IF NOT EXISTS kb_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic TEXT NOT NULL,
    content TEXT NOT NULL,
    embedding BLOB, -- Vector embedding
    metadata TEXT, -- JSON (source, rules, examples)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Projects table (collections of songs)
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    song_ids TEXT, -- JSON array of song IDs
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Collaboration system tables

-- Collaboration sessions
CREATE TABLE IF NOT EXISTS collaboration_sessions (
    id TEXT PRIMARY KEY,
    document_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    created_by TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    settings TEXT, -- JSON
    max_participants INTEGER DEFAULT 10,
    FOREIGN KEY (document_id) REFERENCES songs(id) ON DELETE CASCADE
);

-- Session participants
CREATE TABLE IF NOT EXISTS session_participants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    username TEXT NOT NULL,
    role TEXT NOT NULL, -- owner, editor, viewer
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    permissions TEXT, -- JSON
    cursor_line INTEGER DEFAULT 0,
    cursor_column INTEGER DEFAULT 0,
    FOREIGN KEY (session_id) REFERENCES collaboration_sessions(id) ON DELETE CASCADE,
    UNIQUE(session_id, user_id)
);

-- Document operations for operational transform
CREATE TABLE IF NOT EXISTS document_operations (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL, -- insert, delete, retain
    position INTEGER NOT NULL,
    content TEXT,
    length INTEGER DEFAULT 0,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    version INTEGER NOT NULL,
    dependencies TEXT, -- JSON array
    FOREIGN KEY (session_id) REFERENCES collaboration_sessions(id) ON DELETE CASCADE
);

-- User presence tracking
CREATE TABLE IF NOT EXISTS user_presence (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    status TEXT NOT NULL, -- online, away, busy, offline
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    current_session TEXT,
    device_info TEXT, -- JSON
    FOREIGN KEY (current_session) REFERENCES collaboration_sessions(id) ON DELETE SET NULL
);

-- Session invitations
CREATE TABLE IF NOT EXISTS session_invitations (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    from_user TEXT NOT NULL,
    to_user TEXT NOT NULL,
    role TEXT NOT NULL,
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    accepted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (session_id) REFERENCES collaboration_sessions(id) ON DELETE CASCADE
);

-- Conflict resolution records
CREATE TABLE IF NOT EXISTS conflict_resolutions (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    conflict_id TEXT NOT NULL,
    strategy TEXT NOT NULL,
    resolved_by TEXT NOT NULL,
    resolved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    final_content TEXT,
    changes TEXT, -- JSON
    metadata TEXT, -- JSON
    FOREIGN KEY (session_id) REFERENCES collaboration_sessions(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_songs_updated ON songs(updated_at);
CREATE INDEX IF NOT EXISTS idx_songs_title ON songs(title);
CREATE INDEX IF NOT EXISTS idx_songs_artist ON songs(artist);
CREATE INDEX IF NOT EXISTS idx_versions_song ON versions(song_id, created_at);
CREATE INDEX IF NOT EXISTS idx_stats_date ON writing_stats(date);
CREATE INDEX IF NOT EXISTS idx_kb_topic ON kb_entries(topic);
CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);

-- Collaboration indexes
CREATE INDEX IF NOT EXISTS idx_collab_sessions_active ON collaboration_sessions(is_active, created_at);
CREATE INDEX IF NOT EXISTS idx_collab_sessions_document ON collaboration_sessions(document_id);
CREATE INDEX IF NOT EXISTS idx_participants_session ON session_participants(session_id, is_active);
CREATE INDEX IF NOT EXISTS idx_participants_user ON session_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_operations_session ON document_operations(session_id, version);
CREATE INDEX IF NOT EXISTS idx_operations_timestamp ON document_operations(timestamp);
CREATE INDEX IF NOT EXISTS idx_presence_status ON user_presence(status, last_seen);
CREATE INDEX IF NOT EXISTS idx_presence_session ON user_presence(current_session);
CREATE INDEX IF NOT EXISTS idx_invitations_to_user ON session_invitations(to_user, accepted, expires_at);
CREATE INDEX IF NOT EXISTS idx_invitations_session ON session_invitations(session_id, expires_at);
`
