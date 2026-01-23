package db

// Schema contains the core database table creation statements for single-user mode.
// This is always executed during database initialization.
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

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_songs_updated ON songs(updated_at);
CREATE INDEX IF NOT EXISTS idx_songs_title ON songs(title);
CREATE INDEX IF NOT EXISTS idx_songs_artist ON songs(artist);
CREATE INDEX IF NOT EXISTS idx_versions_song ON versions(song_id, created_at);
CREATE INDEX IF NOT EXISTS idx_stats_date ON writing_stats(date);
CREATE INDEX IF NOT EXISTS idx_kb_topic ON kb_entries(topic);
CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
`

// CollaborationSchema contains database tables for multi-user collaboration features.
// This is only executed when collaboration is enabled via config.Features.EnableCollaboration.
// These tables are a FUTURE FEATURE and are not used in single-user mode.
const CollaborationSchema = `
-- Collaboration system tables (FUTURE FEATURE)
-- These tables are only created when collaboration is enabled.

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

// SyncSchema contains database tables for PWA synchronization.
// This stores captured ideas from companion devices.
const SyncSchema = `
-- Captured ideas from PWA companion
CREATE TABLE IF NOT EXISTS captured_ideas (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL, -- text, voice_memo, photo, tempo
    content TEXT,
    media_path TEXT,
    bpm INTEGER,
    device_id TEXT NOT NULL,
    song_id INTEGER, -- Optional association with a song
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    synced_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE SET NULL
);

-- Paired devices
CREATE TABLE IF NOT EXISTS paired_devices (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    paired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sync indexes
CREATE INDEX IF NOT EXISTS idx_ideas_device ON captured_ideas(device_id, synced_at);
CREATE INDEX IF NOT EXISTS idx_ideas_song ON captured_ideas(song_id);
CREATE INDEX IF NOT EXISTS idx_ideas_type ON captured_ideas(type, created_at);
CREATE INDEX IF NOT EXISTS idx_devices_last_seen ON paired_devices(last_seen);
`

// MuseAgentSchema contains database tables for the AI companion agent.
// This stores conversation history, user preferences, and episodic memories.
const MuseAgentSchema = `
-- Muse agent episodic memory (time-anchored events)
CREATE TABLE IF NOT EXISTS muse_episodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    event_type TEXT NOT NULL,  -- 'edit', 'brainstorm', 'suggestion_accepted', 'suggestion_dismissed', 'chat', 'tool_use'
    song_id INTEGER,
    section TEXT,
    content_snippet TEXT,
    outcome TEXT,
    metadata TEXT,  -- JSON for additional context
    FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE SET NULL
);

-- Muse agent user preferences (learned patterns and settings)
CREATE TABLE IF NOT EXISTS muse_preferences (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,  -- JSON value
    confidence REAL DEFAULT 0.5,
    source TEXT,  -- 'explicit', 'inferred', 'default'
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Muse agent conversation history
CREATE TABLE IF NOT EXISTS muse_conversations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    role TEXT NOT NULL,  -- 'user' or 'assistant'
    content TEXT NOT NULL,
    context TEXT,  -- JSON context (song, section, progress state)
    tool_calls TEXT,  -- JSON array of tool calls made
    tokens_used INTEGER DEFAULT 0
);

-- Muse agent session summary (for long-term context)
CREATE TABLE IF NOT EXISTS muse_session_summaries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT UNIQUE NOT NULL,
    started_at DATETIME NOT NULL,
    ended_at DATETIME,
    summary TEXT,  -- AI-generated session summary
    songs_worked_on TEXT,  -- JSON array of song IDs
    key_insights TEXT,  -- JSON array of insights
    suggestions_accepted INTEGER DEFAULT 0,
    suggestions_dismissed INTEGER DEFAULT 0,
    words_written INTEGER DEFAULT 0
);

-- Muse agent indexes
CREATE INDEX IF NOT EXISTS idx_muse_episodes_session ON muse_episodes(session_id);
CREATE INDEX IF NOT EXISTS idx_muse_episodes_song ON muse_episodes(song_id);
CREATE INDEX IF NOT EXISTS idx_muse_episodes_type ON muse_episodes(event_type, timestamp);
CREATE INDEX IF NOT EXISTS idx_muse_conversations_session ON muse_conversations(session_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_muse_sessions_dates ON muse_session_summaries(started_at, ended_at);
`
