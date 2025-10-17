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

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_songs_updated ON songs(updated_at);
CREATE INDEX IF NOT EXISTS idx_songs_title ON songs(title);
CREATE INDEX IF NOT EXISTS idx_songs_artist ON songs(artist);
CREATE INDEX IF NOT EXISTS idx_versions_song ON versions(song_id, created_at);
CREATE INDEX IF NOT EXISTS idx_stats_date ON writing_stats(date);
CREATE INDEX IF NOT EXISTS idx_kb_topic ON kb_entries(topic);
CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
`
