package agent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/kyanite/noise/internal/domain"
	"github.com/kyanite/noise/internal/infra/db"
	"github.com/kyanite/noise/internal/logging"
	"github.com/google/uuid"
)

// MemoryManager handles episodic, semantic, and working memory for the Muse agent
type MemoryManager struct {
	db        *db.DB
	sessionID string
	working   *WorkingMemory
	config    *AgentConfig
	mutex     sync.RWMutex
}

// NewMemoryManager creates a new memory manager
func NewMemoryManager(database *db.DB, config *AgentConfig) *MemoryManager {
	if config == nil {
		config = DefaultAgentConfig()
	}

	return &MemoryManager{
		db:        database,
		sessionID: uuid.New().String(),
		working: &WorkingMemory{
			SessionStart:  time.Now(),
			LastEdit:      time.Now(),
			RecentEdits:   make([]EditEvent, 0, 100),
			ProgressState: StateStarting,
		},
		config: config,
	}
}

// GetSessionID returns the current session ID
func (m *MemoryManager) GetSessionID() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.sessionID
}

// GetWorkingMemory returns a copy of the current working memory
func (m *MemoryManager) GetWorkingMemory() *WorkingMemory {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy to prevent external mutation
	wm := *m.working
	return &wm
}

// SetCurrentSong updates the current song in working memory
func (m *MemoryManager) SetCurrentSong(song *domain.Song) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.working.CurrentSong = song
}

// SetCurrentSection updates the current section in working memory
func (m *MemoryManager) SetCurrentSection(section string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.working.CurrentSection = section
}

// RecordEdit records an edit event to working memory
func (m *MemoryManager) RecordEdit(edit EditEvent) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	edit.Timestamp = time.Now()
	m.working.RecentEdits = append(m.working.RecentEdits, edit)
	m.working.LastEdit = edit.Timestamp
	m.working.WordsWritten += edit.WordsDelta

	// Keep only recent edits (last 100)
	if len(m.working.RecentEdits) > 100 {
		m.working.RecentEdits = m.working.RecentEdits[len(m.working.RecentEdits)-100:]
	}

	// Update progress state based on edits
	m.updateProgressState()
}

// updateProgressState calculates the current progress state based on recent activity
func (m *MemoryManager) updateProgressState() {
	// Must be called with lock held

	if len(m.working.RecentEdits) == 0 {
		m.working.ProgressState = StateStarting
		return
	}

	now := time.Now()
	timeSinceLastEdit := now.Sub(m.working.LastEdit)

	// Check for stuck state (long pause)
	if timeSinceLastEdit > m.config.StuckThreshold {
		m.working.ProgressState = StateStuck
		return
	}

	// Analyze recent edits for patterns
	recentEdits := m.working.RecentEdits
	if len(recentEdits) < 5 {
		m.working.ProgressState = StateStarting
		return
	}

	// Look at last 10 edits
	start := len(recentEdits) - 10
	if start < 0 {
		start = 0
	}
	recent := recentEdits[start:]

	var insertCount, deleteCount int
	var totalWordsDelta int
	for _, edit := range recent {
		if edit.Type == "insert" {
			insertCount++
		} else if edit.Type == "delete" {
			deleteCount++
		}
		totalWordsDelta += edit.WordsDelta
	}

	// Determine state based on patterns
	if deleteCount > insertCount {
		m.working.ProgressState = StateRefining
	} else if totalWordsDelta > 20 {
		m.working.ProgressState = StateFlowing
	} else if totalWordsDelta < 5 && timeSinceLastEdit > 30*time.Second {
		m.working.ProgressState = StateReviewing
	} else {
		m.working.ProgressState = StateFlowing
	}
}

// UpdateProgressState allows external update of progress state
func (m *MemoryManager) UpdateProgressState(state ProgressState) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.working.ProgressState = state
}

// RecordEpisode records an episodic event to the database
func (m *MemoryManager) RecordEpisode(event EpisodicEvent) error {
	if m.db == nil {
		return fmt.Errorf("database not available")
	}

	event.SessionID = m.GetSessionID()
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	metadataJSON, err := json.Marshal(event.Metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	query := `
		INSERT INTO muse_episodes (session_id, timestamp, event_type, song_id, section, content_snippet, outcome, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = m.db.Exec(query,
		event.SessionID,
		event.Timestamp,
		string(event.EventType),
		event.SongID,
		event.Section,
		event.ContentSnippet,
		event.Outcome,
		string(metadataJSON),
	)

	return err
}

// GetRecentEpisodes retrieves recent episodic events
func (m *MemoryManager) GetRecentEpisodes(limit int) ([]EpisodicEvent, error) {
	if m.db == nil {
		return nil, fmt.Errorf("database not available")
	}

	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT id, session_id, timestamp, event_type, song_id, section, content_snippet, outcome, metadata
		FROM muse_episodes
		WHERE session_id = ?
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := m.db.Query(query, m.GetSessionID(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []EpisodicEvent
	for rows.Next() {
		var event EpisodicEvent
		var songID sql.NullInt64
		var section, contentSnippet, outcome sql.NullString
		var metadataJSON string

		err := rows.Scan(
			&event.ID,
			&event.SessionID,
			&event.Timestamp,
			&event.EventType,
			&songID,
			&section,
			&contentSnippet,
			&outcome,
			&metadataJSON,
		)
		if err != nil {
			logging.Errorf("Failed to scan episode row: %v", err)
			continue
		}

		if songID.Valid {
			id := int(songID.Int64)
			event.SongID = &id
		}
		event.Section = section.String
		event.ContentSnippet = contentSnippet.String
		event.Outcome = outcome.String

		if metadataJSON != "" {
			if err := json.Unmarshal([]byte(metadataJSON), &event.Metadata); err != nil {
				logging.Errorf("Failed to unmarshal episode metadata: %v", err)
			}
		}

		events = append(events, event)
	}

	return events, nil
}

// GetEpisodesForSong retrieves episodes for a specific song
func (m *MemoryManager) GetEpisodesForSong(songID int, limit int) ([]EpisodicEvent, error) {
	if m.db == nil {
		return nil, fmt.Errorf("database not available")
	}

	if limit <= 0 {
		limit = 20
	}

	query := `
		SELECT id, session_id, timestamp, event_type, song_id, section, content_snippet, outcome, metadata
		FROM muse_episodes
		WHERE song_id = ?
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := m.db.Query(query, songID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []EpisodicEvent
	for rows.Next() {
		var event EpisodicEvent
		var sID sql.NullInt64
		var section, contentSnippet, outcome sql.NullString
		var metadataJSON string

		err := rows.Scan(
			&event.ID,
			&event.SessionID,
			&event.Timestamp,
			&event.EventType,
			&sID,
			&section,
			&contentSnippet,
			&outcome,
			&metadataJSON,
		)
		if err != nil {
			logging.Errorf("Failed to scan episode row for song: %v", err)
			continue
		}

		if sID.Valid {
			id := int(sID.Int64)
			event.SongID = &id
		}
		event.Section = section.String
		event.ContentSnippet = contentSnippet.String
		event.Outcome = outcome.String

		if metadataJSON != "" {
			if err := json.Unmarshal([]byte(metadataJSON), &event.Metadata); err != nil {
				logging.Errorf("Failed to unmarshal episode metadata for song: %v", err)
			}
		}

		events = append(events, event)
	}

	return events, nil
}

// SetPreference stores a user preference
func (m *MemoryManager) SetPreference(key string, value interface{}, confidence float64, source string) error {
	if m.db == nil {
		return fmt.Errorf("database not available")
	}

	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO muse_preferences (key, value, confidence, source, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			confidence = excluded.confidence,
			source = excluded.source,
			updated_at = excluded.updated_at
	`

	_, err = m.db.Exec(query, key, string(valueJSON), confidence, source, time.Now())
	return err
}

// GetPreference retrieves a user preference
func (m *MemoryManager) GetPreference(key string) (*Preference, error) {
	if m.db == nil {
		return nil, fmt.Errorf("database not available")
	}

	query := `
		SELECT key, value, confidence, source, updated_at
		FROM muse_preferences
		WHERE key = ?
	`

	var pref Preference
	var valueJSON string

	err := m.db.QueryRow(query, key).Scan(
		&pref.Key,
		&valueJSON,
		&pref.Confidence,
		&pref.Source,
		&pref.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(valueJSON), &pref.Value); err != nil {
		pref.Value = nil
	}
	return &pref, nil
}

// GetAllPreferences retrieves all user preferences
func (m *MemoryManager) GetAllPreferences() (map[string]*Preference, error) {
	if m.db == nil {
		return nil, fmt.Errorf("database not available")
	}

	query := `
		SELECT key, value, confidence, source, updated_at
		FROM muse_preferences
	`

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prefs := make(map[string]*Preference)
	for rows.Next() {
		var pref Preference
		var valueJSON string

		err := rows.Scan(
			&pref.Key,
			&valueJSON,
			&pref.Confidence,
			&pref.Source,
			&pref.UpdatedAt,
		)
		if err != nil {
			logging.Errorf("Failed to scan preference row: %v", err)
			continue
		}

		if err := json.Unmarshal([]byte(valueJSON), &pref.Value); err != nil {
			logging.Errorf("Failed to unmarshal preference value for key %s: %v", pref.Key, err)
		}
		prefs[pref.Key] = &pref
	}

	return prefs, nil
}

// RecordChatMessage records a chat message to the database
func (m *MemoryManager) RecordChatMessage(msg ChatMessage) error {
	if m.db == nil {
		return fmt.Errorf("database not available")
	}

	msg.SessionID = m.GetSessionID()
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	contextJSON, err := json.Marshal(msg.Context)
	if err != nil {
		logging.Errorf("Failed to marshal chat context: %v", err)
		contextJSON = []byte("{}")
	}
	toolCallsJSON, err := json.Marshal(msg.ToolCalls)
	if err != nil {
		logging.Errorf("Failed to marshal chat tool calls: %v", err)
		toolCallsJSON = []byte("[]")
	}

	query := `
		INSERT INTO muse_conversations (session_id, timestamp, role, content, context, tool_calls)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = m.db.Exec(query,
		msg.SessionID,
		msg.Timestamp,
		msg.Role,
		msg.Content,
		string(contextJSON),
		string(toolCallsJSON),
	)

	return err
}

// GetChatHistory retrieves chat history for the current session
func (m *MemoryManager) GetChatHistory(limit int) ([]ChatMessage, error) {
	if m.db == nil {
		return nil, fmt.Errorf("database not available")
	}

	if limit <= 0 {
		limit = m.config.MaxChatHistory
	}

	query := `
		SELECT id, session_id, timestamp, role, content, context, tool_calls
		FROM muse_conversations
		WHERE session_id = ?
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := m.db.Query(query, m.GetSessionID(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		var contextJSON, toolCallsJSON sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.SessionID,
			&msg.Timestamp,
			&msg.Role,
			&msg.Content,
			&contextJSON,
			&toolCallsJSON,
		)
		if err != nil {
			continue
		}

		if contextJSON.Valid && contextJSON.String != "" {
			if err := json.Unmarshal([]byte(contextJSON.String), &msg.Context); err != nil {
				logging.Errorf("Failed to unmarshal chat context: %v", err)
			}
		}
		if toolCallsJSON.Valid && toolCallsJSON.String != "" {
			if err := json.Unmarshal([]byte(toolCallsJSON.String), &msg.ToolCalls); err != nil {
				logging.Errorf("Failed to unmarshal chat tool calls: %v", err)
			}
		}

		messages = append(messages, msg)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// StartSession initializes a new session
func (m *MemoryManager) StartSession() error {
	m.mutex.Lock()
	m.sessionID = uuid.New().String()
	m.working = &WorkingMemory{
		SessionStart:  time.Now(),
		LastEdit:      time.Now(),
		RecentEdits:   make([]EditEvent, 0, 100),
		ProgressState: StateStarting,
	}
	m.mutex.Unlock()

	// Record session start event
	return m.RecordEpisode(EpisodicEvent{
		EventType: EventTypeSessionStart,
		Metadata: map[string]string{
			"started_at": time.Now().Format(time.RFC3339),
		},
	})
}

// EndSession saves session summary and cleans up
func (m *MemoryManager) EndSession() error {
	wm := m.GetWorkingMemory()

	// Record session end event
	err := m.RecordEpisode(EpisodicEvent{
		EventType: EventTypeSessionEnd,
		Metadata: map[string]string{
			"ended_at":      time.Now().Format(time.RFC3339),
			"words_written": fmt.Sprintf("%d", wm.WordsWritten),
			"duration":      time.Since(wm.SessionStart).String(),
		},
	})
	if err != nil {
		return err
	}

	// Save session summary
	return m.saveSessionSummary()
}

// saveSessionSummary creates a summary of the current session
func (m *MemoryManager) saveSessionSummary() error {
	if m.db == nil {
		return fmt.Errorf("database not available")
	}

	wm := m.GetWorkingMemory()
	sessionID := m.GetSessionID()

	// Get songs worked on
	var songsWorkedOn []int
	if wm.CurrentSong != nil {
		songsWorkedOn = append(songsWorkedOn, wm.CurrentSong.ID)
	}
	songsJSON, err := json.Marshal(songsWorkedOn)
	if err != nil {
		logging.Errorf("Failed to marshal songs worked on: %v", err)
		songsJSON = []byte("[]")
	}

	query := `
		INSERT INTO muse_session_summaries (session_id, started_at, ended_at, songs_worked_on, words_written)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(session_id) DO UPDATE SET
			ended_at = excluded.ended_at,
			songs_worked_on = excluded.songs_worked_on,
			words_written = excluded.words_written
	`

	_, err = m.db.Exec(query,
		sessionID,
		wm.SessionStart,
		time.Now(),
		string(songsJSON),
		wm.WordsWritten,
	)

	return err
}

// ClearMemory clears all memory for the current user (privacy feature)
func (m *MemoryManager) ClearMemory(clearEpisodes, clearPreferences, clearConversations bool) error {
	if m.db == nil {
		return fmt.Errorf("database not available")
	}

	if clearEpisodes {
		if _, err := m.db.Exec("DELETE FROM muse_episodes"); err != nil {
			return err
		}
		if _, err := m.db.Exec("DELETE FROM muse_session_summaries"); err != nil {
			return err
		}
	}

	if clearPreferences {
		if _, err := m.db.Exec("DELETE FROM muse_preferences"); err != nil {
			return err
		}
	}

	if clearConversations {
		if _, err := m.db.Exec("DELETE FROM muse_conversations"); err != nil {
			return err
		}
	}

	return nil
}

// GetMemoryStats returns statistics about stored memory
func (m *MemoryManager) GetMemoryStats() (map[string]interface{}, error) {
	if m.db == nil {
		return nil, fmt.Errorf("database not available")
	}

	stats := make(map[string]interface{})

	// Count episodes
	var episodeCount int
	if err := m.db.QueryRow("SELECT COUNT(*) FROM muse_episodes").Scan(&episodeCount); err != nil {
		logging.Errorf("Failed to get episode count: %v", err)
	}
	stats["episode_count"] = episodeCount

	// Count preferences
	var prefCount int
	if err := m.db.QueryRow("SELECT COUNT(*) FROM muse_preferences").Scan(&prefCount); err != nil {
		logging.Errorf("Failed to get preference count: %v", err)
	}
	stats["preference_count"] = prefCount

	// Count conversations
	var convCount int
	if err := m.db.QueryRow("SELECT COUNT(*) FROM muse_conversations").Scan(&convCount); err != nil {
		logging.Errorf("Failed to get conversation count: %v", err)
	}
	stats["conversation_count"] = convCount

	// Count sessions
	var sessionCount int
	if err := m.db.QueryRow("SELECT COUNT(*) FROM muse_session_summaries").Scan(&sessionCount); err != nil {
		logging.Errorf("Failed to get session count: %v", err)
	}
	stats["session_count"] = sessionCount

	// Current session info
	wm := m.GetWorkingMemory()
	stats["current_session"] = map[string]interface{}{
		"session_id":     m.GetSessionID(),
		"started_at":     wm.SessionStart,
		"words_written":  wm.WordsWritten,
		"progress_state": wm.ProgressState.String(),
		"edit_count":     len(wm.RecentEdits),
	}

	return stats, nil
}
