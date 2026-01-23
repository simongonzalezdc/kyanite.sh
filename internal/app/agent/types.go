// Package agent provides the Muse AI companion agent for noise.sh.
// It includes memory management, context observation, proactive suggestions,
// chat interface, tool integration, and learning capabilities.
package agent

import (
	"time"

	"github.com/Kyanite/noise/internal/domain"
)

// EventType represents the type of episodic event
type EventType string

const (
	EventTypeEdit                EventType = "edit"
	EventTypeBrainstorm          EventType = "brainstorm"
	EventTypeSuggestionAccepted  EventType = "suggestion_accepted"
	EventTypeSuggestionDismissed EventType = "suggestion_dismissed"
	EventTypeChat                EventType = "chat"
	EventTypeToolUse             EventType = "tool_use"
	EventTypeSessionStart        EventType = "session_start"
	EventTypeSessionEnd          EventType = "session_end"
)

// ProgressState represents the user's current writing progress state
type ProgressState int

const (
	StateStarting  ProgressState = iota // New song, blank page
	StateFlowing                        // Good velocity, consistent edits
	StateStuck                          // Long pause, deletions > additions
	StateRefining                       // Small tweaks, low velocity
	StateReviewing                      // Reading, no edits
)

// String returns the string representation of ProgressState
func (s ProgressState) String() string {
	switch s {
	case StateStarting:
		return "starting"
	case StateFlowing:
		return "flowing"
	case StateStuck:
		return "stuck"
	case StateRefining:
		return "refining"
	case StateReviewing:
		return "reviewing"
	default:
		return "unknown"
	}
}

// EpisodicEvent represents a time-anchored memory event
type EpisodicEvent struct {
	ID             int64             `json:"id"`
	SessionID      string            `json:"session_id"`
	Timestamp      time.Time         `json:"timestamp"`
	EventType      EventType         `json:"event_type"`
	SongID         *int              `json:"song_id,omitempty"`
	Section        string            `json:"section,omitempty"`
	ContentSnippet string            `json:"content_snippet,omitempty"`
	Outcome        string            `json:"outcome,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// Preference represents a learned user preference
type Preference struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	Confidence float64     `json:"confidence"`
	Source     string      `json:"source"` // "explicit", "inferred", "default"
	UpdatedAt  time.Time   `json:"updated_at"`
}

// ChatMessage represents a message in the conversation
type ChatMessage struct {
	ID        int64             `json:"id"`
	SessionID string            `json:"session_id"`
	Timestamp time.Time         `json:"timestamp"`
	Role      string            `json:"role"` // "user" or "assistant"
	Content   string            `json:"content"`
	Context   map[string]string `json:"context,omitempty"`
	ToolCalls []ToolCall        `json:"tool_calls,omitempty"`
}

// ToolCall represents a tool invocation by the agent
type ToolCall struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
	Result string            `json:"result,omitempty"`
}

// SessionSummary represents a summary of a writing session
type SessionSummary struct {
	ID                    int64     `json:"id"`
	SessionID             string    `json:"session_id"`
	StartedAt             time.Time `json:"started_at"`
	EndedAt               time.Time `json:"ended_at,omitempty"`
	Summary               string    `json:"summary,omitempty"`
	SongsWorkedOn         []int     `json:"songs_worked_on,omitempty"`
	KeyInsights           []string  `json:"key_insights,omitempty"`
	SuggestionsAccepted   int       `json:"suggestions_accepted"`
	SuggestionsDismissed  int       `json:"suggestions_dismissed"`
	WordsWritten          int       `json:"words_written"`
}

// WorkingMemory holds the current session state
type WorkingMemory struct {
	CurrentSong    *domain.Song  `json:"current_song,omitempty"`
	CurrentSection string        `json:"current_section,omitempty"`
	RecentEdits    []EditEvent   `json:"recent_edits,omitempty"`
	SessionStart   time.Time     `json:"session_start"`
	LastEdit       time.Time     `json:"last_edit"`
	WordsWritten   int           `json:"words_written"`
	ProgressState  ProgressState `json:"progress_state"`
}

// EditEvent represents a single edit action
type EditEvent struct {
	Timestamp   time.Time `json:"timestamp"`
	Type        string    `json:"type"` // "insert", "delete", "replace"
	Position    int       `json:"position"`
	Length      int       `json:"length"`
	Content     string    `json:"content,omitempty"`
	WordsDelta  int       `json:"words_delta"`
}

// SuggestionType represents the type of suggestion
type SuggestionType int

const (
	SuggestNextLine SuggestionType = iota
	SuggestChord
	SuggestBreak
	SuggestReference
	SuggestStructure
	SuggestRhyme
	SuggestGoal
)

// String returns the string representation of SuggestionType
func (s SuggestionType) String() string {
	switch s {
	case SuggestNextLine:
		return "next_line"
	case SuggestChord:
		return "chord"
	case SuggestBreak:
		return "break"
	case SuggestReference:
		return "reference"
	case SuggestStructure:
		return "structure"
	case SuggestRhyme:
		return "rhyme"
	case SuggestGoal:
		return "goal"
	default:
		return "unknown"
	}
}

// Suggestion represents a proactive suggestion from the agent
type Suggestion struct {
	ID         string            `json:"id"`
	Type       SuggestionType    `json:"type"`
	Content    string            `json:"content"`
	Confidence float64           `json:"confidence"`
	Reason     string            `json:"reason"`
	Actions    []SuggestionAction `json:"actions,omitempty"`
	ExpiresAt  time.Time         `json:"expires_at"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// SuggestionAction represents an action the user can take on a suggestion
type SuggestionAction struct {
	Label    string `json:"label"`
	Key      string `json:"key"` // keyboard shortcut
	ActionID string `json:"action_id"`
}

// UserPreferences holds all learned user preferences
type UserPreferences struct {
	PreferredGenres      []string `json:"preferred_genres,omitempty"`
	TypicalStructure     string   `json:"typical_structure,omitempty"` // "ABABCB", "AABA", etc.
	LyricStyle           string   `json:"lyric_style,omitempty"`       // "narrative", "abstract", "conversational"
	ChordPreferences     []string `json:"chord_preferences,omitempty"`
	WorkingHours         []int    `json:"working_hours,omitempty"` // Hours of day user typically works
	SuggestionFrequency  string   `json:"suggestion_frequency,omitempty"` // "minimal", "moderate", "eager"
	PreferredAIVoice     string   `json:"preferred_ai_voice,omitempty"`   // "formal", "casual", "encouraging"
}

// AgentConfig holds configuration for the Muse agent
type AgentConfig struct {
	// Memory settings
	MaxEpisodicEvents     int           `json:"max_episodic_events"`
	SessionTimeout        time.Duration `json:"session_timeout"`
	
	// Suggestion settings
	SuggestionCooldown    time.Duration `json:"suggestion_cooldown"`
	MinConfidence         float64       `json:"min_confidence"`
	StuckThreshold        time.Duration `json:"stuck_threshold"`
	
	// Chat settings
	MaxChatHistory        int           `json:"max_chat_history"`
	ContextWindowSize     int           `json:"context_window_size"`
	
	// Learning settings
	LearningEnabled       bool          `json:"learning_enabled"`
	PreferenceUpdateRate  float64       `json:"preference_update_rate"`
}

// DefaultAgentConfig returns the default agent configuration
func DefaultAgentConfig() *AgentConfig {
	return &AgentConfig{
		MaxEpisodicEvents:    1000,
		SessionTimeout:       2 * time.Hour,
		SuggestionCooldown:   5 * time.Minute,
		MinConfidence:        0.6,
		StuckThreshold:       3 * time.Minute,
		MaxChatHistory:       100,
		ContextWindowSize:    10,
		LearningEnabled:      true,
		PreferenceUpdateRate: 0.1,
	}
}
